package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Listing struct {
	Id string `json:"id"`
}

type Error struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func getAzError(err error) *azcore.ResponseError {
	var responseErr *azcore.ResponseError
	errors.As(err, &responseErr)
	return responseErr
}

func getContainer() (*azcosmos.ContainerClient, error) {
	endpoint, ok := os.LookupEnv("DB_URI")
	if !ok {
		log.Fatal("DB_URI could not be found")
		return nil, errors.New("DB_URI could not be found")
	}

	key, ok := os.LookupEnv("DB_PRIMARY_KEY")
	if !ok {
		log.Fatal("DB_PRIMARY_KEY could not be found")
		return nil, errors.New("DB_PRIMARY_KEY could not be found")
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		log.Fatal("DB_NAME could not be found")
		return nil, errors.New("DB_NAME could not be found")
	}

	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		log.Fatal("Failed to create credentials from DB_PRIMARY_KEY")
		return nil, err
	}

	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		log.Fatal("Failed to create client")
		return nil, err
	}

	container, err := client.NewContainer(dbName, "listings")
	if err != nil {
		log.Fatal("Failed to get container Listings")
		return nil, err
	}
	return container, nil
}

func getListings(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")

	container, err := getContainer()
	if err != nil {
		render.JSON(w, r, Error{Message: err.Error(), StatusCode: http.StatusInternalServerError})
		return
	}

	query := "SELECT * FROM c"
	partitionKey := azcosmos.NewPartitionKeyString(strings.ToLower(city))
	pager := container.NewQueryItemsPager(query, partitionKey, nil)

	listings := []Listing{}

	for pager.More() {
		response, err := pager.NextPage(context.Background())
		if err != nil {
			var azError *azcore.ResponseError
			errors.As(err, &azError)
			log.Printf("Failed to get next result page: %s\n", azError.ErrorCode)
			render.JSON(w, r, Error{Message: err.Error(), StatusCode: http.StatusInternalServerError})
			return
		}

		for _, bytes := range response.Items {
			listing := Listing{}
			if err := json.Unmarshal(bytes, &listing); err != nil {
				log.Printf("An error occurred trying to parse the response json: %s", err.Error())
				render.JSON(w, r, Error{Message: err.Error(), StatusCode: http.StatusInternalServerError})
				return
			}
			listings = append(listings, listing)
		}
	}

	d, _ := json.Marshal(listings)
	fmt.Fprintln(w, string(d))
}

func getListingsById(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")
	listingId := chi.URLParam(r, "listingId")

	container, err := getContainer()
	if err != nil {
		render.JSON(w, r, Error{Message: err.Error(), StatusCode: http.StatusInternalServerError})
		return
	}

	partitionKey := azcosmos.NewPartitionKeyString(strings.ToLower(city))
	response, err := container.ReadItem(context.Background(), partitionKey, listingId, nil)
	if err != nil {
		var azError *azcore.ResponseError
		errors.As(err, &azError)

		if azError.StatusCode == http.StatusNotFound {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, Error{
				Message:    fmt.Sprintf("Listing with id %s could not be found", listingId),
				StatusCode: http.StatusNotFound,
			})
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Printf("Failed to read item: %s\n", azError.ErrorCode)
		render.JSON(w, r, Error{Message: err.Error(), StatusCode: http.StatusInternalServerError})
		return
	}

	if response.RawResponse.StatusCode == 200 {
		listing := Listing{}
		err := json.Unmarshal(response.Value, &listing)
		if err != nil {
			log.Printf("Failed to unmarshal result: %s\n", err.Error())
			render.JSON(w, r, Error{Message: err.Error(), StatusCode: http.StatusInternalServerError})
			return
		}

		render.JSON(w, r, listing)
	} else {
		log.Printf("Item could not successfully be fetched")
		render.JSON(w, r, Error{
			Message:    "Item could not successfully be fetched",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
}

func main() {
	port := "8080"
	log.Printf("About to listen on %s. Go to http://127.0.0.1:%s/", port, port)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/{city}/listings", getListings)
	r.Get("/{city}/listings/{listingId}", getListingsById)

	http.ListenAndServe(":"+port, r)
}
