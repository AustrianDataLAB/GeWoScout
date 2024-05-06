package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/AustrianDataLAB/GeWoScout/backend/cosmos"
	"github.com/AustrianDataLAB/GeWoScout/backend/models"

	"github.com/go-chi/render"
)

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
func GetListings(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://enmwjn0cdi37l.x.pipedream.net") // TODO addded just for debugging

	req, err := models.InvokeRequestFromBody(r.Body)
	if err != nil {
		log.Printf("Failed to read invoke request body: %s\n", err.Error())
		render.JSON(w, r, models.NewInvokeResponse(
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
		render.JSON(w, r, models.NewInvokeResponse(
			http.StatusBadRequest,
			models.Error{Message: "City param was invalid or empty", StatusCode: http.StatusBadRequest},
		))
		return
	}

	container, err := cosmos.GetContainer()
	if err != nil {
		render.JSON(w, r, models.NewInvokeResponse(
			http.StatusInternalServerError,
			models.Error{Message: err.Error(), StatusCode: http.StatusInternalServerError},
		))
		return
	}

	pager := cosmos.GetQueryItemsPager(container, city, &req.Data.Req.Query)
	listings := []models.Listing{}
	var continuationToken *string

	for pager.More() {
		response, err := pager.NextPage(context.Background())
		if err != nil {
			log.Printf("Failed to get next result page: %s\n", err.Error())

			render.JSON(w, r, models.NewInvokeResponse(
				http.StatusBadRequest,
				models.Error{Message: err.Error(), StatusCode: http.StatusBadRequest},
			))
			return
		}

		for _, bytes := range response.Items {
			listing := models.Listing{}
			if err := json.Unmarshal(bytes, &listing); err != nil {
				log.Printf("An error occurred trying to parse the response json: %s", err.Error())

				render.JSON(w, r, models.NewInvokeResponse(
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

	render.JSON(w, r, models.NewInvokeResponse(http.StatusOK, result))
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
func GetListingById(w http.ResponseWriter, r *http.Request) {
	injectedData := models.CosmosBindingInput{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&injectedData); err != nil {
		log.Printf("Error trying to unmarshal injected data: %s\n", err.Error())
		return
	}

	if injectedData.Data.Documents == "null" {
		render.JSON(w, r, models.NewInvokeResponse(
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

		render.JSON(w, r, models.NewInvokeResponse(
			http.StatusBadRequest,
			models.Error{
				Message:    fmt.Sprintf("Listing with id %s could not be found in city %s", injectedData.Metadata.ID, injectedData.Metadata.City),
				StatusCode: http.StatusBadRequest,
			},
		))
		return
	}

	render.JSON(w, r, models.NewInvokeResponse(http.StatusOK, listing))
}
