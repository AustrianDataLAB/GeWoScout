package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/go-chi/render"
	"github.com/xeipuuv/gojsonschema"

	"github.com/AustrianDataLAB/GeWoScout/backend/cosmos"
)

const (
	XMSClientPrincipalHeaderName = "X-MS-CLIENT-PRINCIPAL"
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
// @Param title query string false "Listing title to search for"
// @Param housingCooperative query string false "Name of the 'Genossenschaft'"
// @Param projectId query string false "Project ID for which to return listings"
// @Param postalCode query string false "Postal code(s) within which to look for listings"
// @Param roomCount query integer false "Exact room count to search for"
// @Param minRoomCount query integer false "Minimum number of rooms"
// @Param maxRoomCount query integer false "Maximum number of rooms"
// @Param minSqm query integer false "Minimum number of square meters"
// @Param maxSqm query integer false "Maximum number of square meters"
// @Param availableFrom query string false "Date from which the listing has to be available (latest date)"
// @Param minYearBuilt query integer false "Oldest allowed construction year"
// @Param maxYearBuilt query integer false "Most recent allowed construction year"
// @Param minHwgEnergyClass query string false "Worst acceptable HWG energy class" Enums(A++, A+, A, B, C, D, E, F)
// @Param minFgeeEnergyClass query string false "Worst acceptable fgEE energy class" Enums(A++, A+, A, B, C, D, E, F)
// @Param listingType query string false "Type of listing" Enums(rent, sale, both)
// @Param minRentPricePerMonth query integer false "Minimum rent per month"
// @Param maxRentPricePerMonth query integer false "Maximum rent per month"
// @Param minCooperativeShare query integer false "Minimum cooperative share"
// @Param maxCooperativeShare query integer false "Maximum cooperative share"
// @Param minSalePrice query integer false "Minimum sale price"
// @Param maxSalePrice query integer false "Maximum sale price"
// @Param sortBy query string false "Field to sort by"
// @Param sortType query string false "Whether to search ascending or descending" Enums(ASC, DESC)
// @Success 200 {object} models.GetListingsResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /cities/{city}/listings [get]
func (h *Handler) GetListings(w http.ResponseWriter, r *http.Request) {
	req, err := models.UnmarshalAndValidate[models.InvokeRequest[models.ListingsQuery]](r.Body)
	if err != nil {
		// Error is returned to the user here because the validation errors
		// return information about which fields were invalid.
		errMsg := fmt.Sprintf("Failed to read invoke request body: %s\n", err.Error())
		render.JSON(w, r, models.NewHttpInvokeResponse(
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
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusBadRequest,
			models.Error{Message: errMsg},
			[]string{errMsg},
		))
		return
	}

	pager := cosmos.GetQueryItemsPager(h.GetListingsByCityContainerClient(), city, &req.Data.Req.Query)

	maxNumListings := cosmos.DEFAULT_PAGE_SIZE
	if req.Data.Req.Query.PageSize != nil {
		maxNumListings = *req.Data.Req.Query.PageSize
	}

	listings := make([]models.Listing, 0, maxNumListings)
	var continuationToken *string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			render.JSON(w, r, models.NewHttpInvokeResponse(
				http.StatusBadRequest,
				models.Error{Message: "Failed to get listings"},
				[]string{fmt.Sprintf("Failed to get next result page: %s\n", err.Error())},
			))
			return
		}

		for _, bytes := range response.Items {
			listing := models.Listing{}
			if err := json.Unmarshal(bytes, &listing); err != nil {
				render.JSON(w, r, models.NewHttpInvokeResponse(
					http.StatusBadRequest,
					models.Error{Message: "Failed to get listings"},
					[]string{fmt.Sprintf("Failed to parse listings: %s\n", err.Error())},
				))
				return
			}
			listings = append(listings, listing)
		}

		continuationToken = response.ContinuationToken
	}

	result := models.GetListingsResponse{
		Results:           listings,
		ContinuationToken: continuationToken,
	}

	render.JSON(w, r, models.NewHttpInvokeResponse(http.StatusOK, result, nil))
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
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusInternalServerError,
			models.Error{
				Message: "Failed to get listing",
			},
			[]string{fmt.Sprintf("Failed to unmarshal injected data: %s\n", err.Error())},
		))
		return
	}

	if injectedData.Data.Documents == "null" {
		errMsg := fmt.Sprintf("Listing with id %s could not be found in city %s", injectedData.Metadata.ID, injectedData.Metadata.City)
		render.JSON(w, r, models.NewHttpInvokeResponse(
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
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusBadRequest,
			models.Error{
				Message: "Failed to get listing",
			},
			[]string{fmt.Sprintf("Failed to unmarshal injected listing: %s\n", err.Error())},
		))
		return
	}

	render.JSON(w, r, models.NewHttpInvokeResponse(http.StatusOK, listing, nil))
}

// UpdateUserPrefs Handler function for /updateUserPrefs, which sets a user's
// notification preferences by either creating or updating them in the DB.
// Neither the city nor the id are guaranteed to exist at this point.
// @Summary Set (or update) a user's notification preferences
// @Description Update notification preferences of the user with the given id
// @Tags userPreferences
// @Accept json
// @Produce json
// @Param city path string true "The city assigned to the user"
// @Param id path string true "The id of the user"
// @Success 200 {object} models.NotificationSettings "Successfully updates notification settings"
// @Failure 404 {object} models.Error "Notification settings could not be updated"
// @Failure 400 {object} models.Error "Bad request"
// @Router /cities/{city}/users/{id}/preferences [put]
func (h *Handler) UpdateUserPrefs(w http.ResponseWriter, r *http.Request) {
	req, err := models.UnmarshalAndValidate[models.InvokeRequest[interface{}]](r.Body)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to read invoke request body: %s\n", err.Error())
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusBadRequest,
			models.Error{Message: errMsg},
			[]string{errMsg},
		))
		return
	}

	if len(req.Data.Req.Body) == 0 {
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusBadRequest,
			models.Error{Message: "Missing body in request data"},
			[]string{"Missing body in request data"},
		))
		return
	}

	_, err = models.UnmarshalAndValidate[models.NotificationSettings](io.NopCloser(strings.NewReader(req.Data.Req.Body)))

	if err != nil {
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusBadRequest,
			models.Error{Message: "Failed to parse request body"},
			[]string{fmt.Sprintf("Failed to parse request body: %s\n", err.Error())},
		))
		return
	}

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
	render.JSON(w, r, models.NewHttpInvokeResponse(http.StatusOK, aliveResponse, nil))
}
