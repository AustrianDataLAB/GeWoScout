package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/go-chi/render"
	"github.com/xeipuuv/gojsonschema"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/AustrianDataLAB/GeWoScout/backend/cosmos"
)

type Handler struct {
	cosmosOnce sync.Once
	// Do NOT access directly
	cosmosClient *azcosmos.Client
	// Do NOT access directly
	gewoscoutDbClient *azcosmos.DatabaseClient
	// Do NOT access directly
	listingsByCityContainerClient *azcosmos.ContainerClient

	ScraperResultSchema *gojsonschema.Schema
}

func (h *Handler) initCosmos() {
	h.cosmosClient, h.gewoscoutDbClient, h.listingsByCityContainerClient = cosmos.InitClients()
}

func (h *Handler) GetCosmosClient() *azcosmos.Client {
	h.cosmosOnce.Do(h.initCosmos)
	return h.cosmosClient
}

func (h *Handler) GetGewoscoutDbClient() *azcosmos.DatabaseClient {
	h.cosmosOnce.Do(h.initCosmos)
	return h.gewoscoutDbClient
}

func (h *Handler) GetListingsByCityContainerClient() *azcosmos.ContainerClient {
	h.cosmosOnce.Do(h.initCosmos)
	return h.listingsByCityContainerClient
}

func NewHandler() *Handler {
	schemaLoader := gojsonschema.NewStringLoader(models.ScraperResultListingSchema)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		log.Fatalf("Failed to create schema: %s", err.Error())
	}

	return &Handler{
		ScraperResultSchema: schema,
	}
}

// GetListings Handler function for /listings, which returns any listings within the
// partition defined by the city path param, which is guaranteed to exist at
// this point.
// @Summary Get listings for a city
// @Description Get listings for a city
// @Tags listings
// @Accept json
// @Produce json
// @Param city path string true "The city for which to get listings"
// @Param continuationToken query string false "The continuation token for pagination"
// @Success 200 {object} models.GetListingsResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /cities/{city}/listings [get]
func (h *Handler) GetListings(w http.ResponseWriter, r *http.Request) {
	req, err := models.InvokeRequestFromBody(r.Body)
	if err != nil {
		log.Printf("Failed to read invoke request body: %s\n", err.Error())
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusBadRequest,
			models.Error{Message: err.Error(), StatusCode: http.StatusBadRequest},
		))
		return
	}
	city := req.Data.Req.Params["city"]

	// Manual check that city is not empty, validation through validator package
	// does not work
	if len(strings.TrimSpace(city)) == 0 {
		log.Println("City param was invalid empty")
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusBadRequest,
			models.Error{Message: "City param was invalid or empty", StatusCode: http.StatusBadRequest},
		))
		return
	}

	pager := cosmos.GetQueryItemsPager(h.GetListingsByCityContainerClient(), city, &req.Data.Req.Query)
	var listings = make([]models.Listing, 0, 30)
	var continuationToken *string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			log.Printf("Failed to get next result page: %s\n", err.Error())

			render.JSON(w, r, models.NewHttpInvokeResponse(
				http.StatusBadRequest,
				models.Error{Message: err.Error(), StatusCode: http.StatusBadRequest},
			))
			return
		}

		for _, bytes := range response.Items {
			listing := models.Listing{}
			if err := json.Unmarshal(bytes, &listing); err != nil {
				log.Printf("An error occurred trying to parse the response json: %s", err.Error())

				render.JSON(w, r, models.NewHttpInvokeResponse(
					http.StatusBadRequest,
					models.Error{Message: err.Error(), StatusCode: http.StatusBadRequest},
				))
				return
			}
			listings = append(listings, listing)
		}

		continuationToken = response.ContinuationToken
		break
	}

	result := models.GetListingsResponse{
		Results:           listings,
		ContinuationToken: continuationToken,
	}

	render.JSON(w, r, models.NewHttpInvokeResponse(http.StatusOK, result))
}

// GetListingById Handler function for /listingById, which returns a listing by its id
// and city. Neither the city nor the id are guaranteed to exist at this point.
// @Summary Get a listing by its id
// @Description Get a listing by its id
// @Tags listings
// @Accept json
// @Produce json
// @Param city path string true "The city for which to get the listing"
// @Param id path string true "The id of the listing"
// @Success 200 {object} models.Listing "Successfully found listing"
// @Failure 404 {object} models.Error "Listing not found"
// @Failure 400 {object} models.Error "Bad request"
// @Router /cities/{city}/listings/{id} [get]
func (h *Handler) GetListingById(w http.ResponseWriter, r *http.Request) {
	injectedData := models.CosmosBindingInput{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&injectedData); err != nil {
		log.Printf("Error trying to unmarshal injected data: %s\n", err.Error())
		return
	}

	if injectedData.Data.Documents == "null" {
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusNotFound,
			models.Error{
				Message:    fmt.Sprintf("Listing with id %s could not be found in city %s", injectedData.Metadata.ID, injectedData.Metadata.City),
				StatusCode: http.StatusNotFound,
			},
		))
		return
	}

	input, _ := strconv.Unquote(injectedData.Data.Documents)

	listing := models.Listing{}
	if err := json.Unmarshal([]byte(input), &listing); err != nil {
		log.Printf("Error trying to unmarshal injected listing: %s\n", err.Error())

		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusBadRequest,
			models.Error{
				Message:    fmt.Sprintf("Listing with id %s could not be found in city %s", injectedData.Metadata.ID, injectedData.Metadata.City),
				StatusCode: http.StatusBadRequest,
			},
		))
		return
	}

	render.JSON(w, r, models.NewHttpInvokeResponse(http.StatusOK, listing))
}

// HandleHealth Handler function for /health, which returns a simple alive response
// @Summary Health check
// @Description Health check
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} models.HealthResponse "Alive"
// @Router /health [get]
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	aliveResponse := models.HealthResponse{
		Status: "ok",
	}
	render.JSON(w, r, models.NewHttpInvokeResponse(http.StatusOK, aliveResponse))
}
