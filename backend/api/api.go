package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/AustrianDataLAB/GeWoScout/backend/cosmos"
	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Handler function for /listings, which returns any listings within the
// partition defined by the city path param, which is guaranteed to exist at
// this point.
func GetListings(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")

	container, err := cosmos.GetContainer()
	if err != nil {
		render.JSON(w, r, models.Error{Message: err.Error(), StatusCode: http.StatusInternalServerError})
		return
	}

	query := "SELECT * FROM c WHERE c._partitionKey = @city"
	partitionKey := azcosmos.NewPartitionKeyString(strings.ToLower(city))
	options := azcosmos.QueryOptions{QueryParameters: []azcosmos.QueryParameter{{Name: "@city", Value: city}}}
	pager := container.NewQueryItemsPager(query, partitionKey, &options)

	listings := []models.Listing{}

	for pager.More() {
		response, err := pager.NextPage(context.Background())
		if err != nil {
			var azError *azcore.ResponseError
			errors.As(err, &azError)
			log.Printf("Failed to get next result page: %s\n", azError.ErrorCode)
			render.JSON(w, r, models.Error{Message: err.Error(), StatusCode: http.StatusInternalServerError})
			return
		}

		for _, bytes := range response.Items {
			listing := models.Listing{}
			if err := json.Unmarshal(bytes, &listing); err != nil {
				log.Printf("An error occurred trying to parse the response json: %s", err.Error())
				render.JSON(w, r, models.Error{Message: err.Error(), StatusCode: http.StatusInternalServerError})
				return
			}
			listings = append(listings, listing)
		}
	}

	result := make(map[string]interface{})

	result["results"] = listings
	result["continuationToken"] = nil // TODO get ct here

	render.JSON(w, r, result)
}

func GetListingById(w http.ResponseWriter, r *http.Request) {
	injectedData := models.CosmosBindingInput{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&injectedData); err != nil {
		log.Printf("Error trying to unmarshal injected data: %s\n", err.Error())
		render.JSON(w, r,
			models.InvokeResponse{
				Outputs: map[string]interface{}{
					"statusCode": http.StatusBadRequest,
				},
				Logs:        []string{},
				ReturnValue: models.Error{Message: err.Error(), StatusCode: http.StatusBadRequest},
			},
		)
		return
	}

	if injectedData.Data.Documents == "null" {
		render.JSON(w, r, models.InvokeResponse{
			Outputs: map[string]interface{}{
				"statusCode": http.StatusNotFound,
			},
			Logs: []string{},
			ReturnValue: models.Error{
				Message:    fmt.Sprintf("Listing with id %s could not be found in city %s", injectedData.Metadata.ID, injectedData.Metadata.City),
				StatusCode: http.StatusNotFound,
			},
		})
		return
	}

	input, _ := strconv.Unquote(injectedData.Data.Documents)

	listing := models.Listing{}
	if err := json.Unmarshal([]byte(input), &listing); err != nil {
		log.Printf("Error trying to unmarshal injected listing: %s\n", err.Error())
		render.JSON(w, r,
			models.InvokeResponse{
				Outputs: map[string]interface{}{
					"statusCode": http.StatusBadRequest,
				},
				Logs:        []string{},
				ReturnValue: models.Error{Message: err.Error(), StatusCode: http.StatusBadRequest},
			},
		)
		return
	}

	outputs := make(map[string]interface{})
	outputs["statusCode"] = http.StatusOK
	invokeResponse := models.InvokeResponse{Outputs: outputs, Logs: []string{}, ReturnValue: listing}
	render.JSON(w, r, invokeResponse)
}
