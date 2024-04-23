package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/AustrianDataLAB/GeWoScout/backend/cosmos"
	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"

	"github.com/go-chi/render"
)

// Handler function for /listings, which returns any listings within the
// partition defined by the city path param, which is guaranteed to exist at
// this point.
func GetListings(w http.ResponseWriter, r *http.Request) {
	req, err := models.InvokeRequestFromBody(r.Body)
	if err != nil {
		render.JSON(w, r, models.NewInvokeResponse(
			http.StatusBadRequest,
			models.Error{Message: err.Error(), StatusCode: http.StatusBadRequest},
		))
	}
	city := req.Data.Req.Params["city"]

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

	for pager.More() {
		response, err := pager.NextPage(context.Background())
		if err != nil {
			var azError *azcore.ResponseError
			errors.As(err, &azError)
			log.Printf("Failed to get next result page: %s\n", azError.ErrorCode)

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
	}

	result := make(map[string]interface{})

	result["results"] = listings
	if req.Data.Req.Query.ContinuationToken != "" {
		result["continuationToken"] = req.Data.Req.Query.ContinuationToken
	} else {
		result["continuationToken"] = nil
	}

	render.JSON(w, r, models.NewInvokeResponse(http.StatusOK, result))
}

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
