package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/go-chi/render"

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
	return &Handler{}
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
		errMsg := fmt.Sprintf("Failed to read invoke request body: %s\n", err.Error())
		render.JSON(w, r, models.NewInvokeResponse(
			http.StatusBadRequest,
			models.Error{Message: errMsg},
			[]string{errMsg},
		))
		return
	}
	city := req.Data.Req.Params["city"]

	// Manual check that city is not empty, validation through validator package
	// does not work
	if len(strings.TrimSpace(city)) == 0 {
		errMsg := "City param was invalid or empty"
		render.JSON(w, r, models.NewInvokeResponse(
			http.StatusBadRequest,
			models.Error{Message: errMsg},
			[]string{errMsg},
		))
		return
	}

	pager := cosmos.GetQueryItemsPager(h.GetListingsByCityContainerClient(), city, &req.Data.Req.Query)
	listings := []models.Listing{}
	var continuationToken *string

	for pager.More() {
		// TODO add timeout
		response, err := pager.NextPage(context.Background())
		if err != nil {
			errMsg := fmt.Sprintf("Failed to get next result page: %s\n", err.Error())
			render.JSON(w, r, models.NewInvokeResponse(
				http.StatusBadRequest,
				models.Error{Message: errMsg},
				[]string{errMsg},
			))
			return
		}

		for _, bytes := range response.Items {
			listing := models.Listing{}
			if err := json.Unmarshal(bytes, &listing); err != nil {
				errMsg := fmt.Sprintf("An error occurred trying to parse the response json: %s", err.Error())
				render.JSON(w, r, models.NewInvokeResponse(
					http.StatusBadRequest,
					models.Error{Message: err.Error()},
					[]string{errMsg},
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

	render.JSON(w, r, models.NewInvokeResponse(http.StatusOK, result, nil))
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
		errMsg := fmt.Sprintf("Error trying to unmarshal injected data: %s\n", err.Error())
		render.JSON(w, r, models.NewInvokeResponse(
			http.StatusInternalServerError,
			models.Error{
				Message: errMsg,
			},
			[]string{errMsg},
		))
		return
	}

	if injectedData.Data.Documents == "null" {
		errMsg := fmt.Sprintf("Listing with id %s could not be found in city %s", injectedData.Metadata.ID, injectedData.Metadata.City)
		render.JSON(w, r, models.NewInvokeResponse(
			http.StatusNotFound,
			models.Error{
				Message: errMsg,
			},
			[]string{errMsg},
		))
		return
	}

	input, _ := strconv.Unquote(injectedData.Data.Documents)

	listing := models.Listing{}
	if err := json.Unmarshal([]byte(input), &listing); err != nil {
		errMsg := fmt.Sprintf("Error trying to unmarshal injected listing: %s\n", err.Error())
		render.JSON(w, r, models.NewInvokeResponse(
			http.StatusBadRequest,
			models.Error{
				Message: errMsg,
			},
			[]string{errMsg},
		))
		return
	}

	render.JSON(w, r, models.NewInvokeResponse(http.StatusOK, listing, nil))
}
