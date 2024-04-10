package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Listing struct {
	ListingId string
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!\n")
}

func getListings(w http.ResponseWriter, r *http.Request) {
	endpoint, ok := os.LookupEnv("DB_URI")
	if !ok {
		log.Fatal("DB_URI could not be found")
	}

	key, ok := os.LookupEnv("DB_PRIMARY_KEY")
	if !ok {
		log.Fatal("DB_PRIMARY_KEY could not be found")
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		log.Fatal("DB_NAME could not be found")
	}

	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		log.Fatal("Failed to create credentials from DB_PRIMARY_KEY")
	}

	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		log.Fatal("Failed to create client")
	}

	container, err := client.NewContainer(dbName, "listings")
	if err != nil {
		log.Fatal("Failed to get container listings")
	}

	query := "SELECT * FROM c"
	partitionKey := azcosmos.NewPartitionKeyString("listingId")
	pager := container.NewQueryItemsPager(query, partitionKey, nil)

	listings := []Listing{}

	for pager.More() {
		response, err := pager.NextPage(context.Background())
		if err != nil {
			var responseErr *azcore.ResponseError
			errors.As(err, &responseErr)
			log.Printf("An error occurred reading from the database: %s", responseErr.Error())
			fmt.Fprintf(w, "Internal error")
			return
		}

		for _, bytes := range response.Items {
			listing := Listing{}
			if err := json.Unmarshal(bytes, &listing); err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			listings = append(listings, listing)
		}
	}

	d, _ := json.Marshal(listings)
	fmt.Fprintf(w, string(d))
}

func main() {
	port := "8080"
	log.Printf("About to listen on %s. Go to http://127.0.0.1:%s/", port, port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/listings", getListings)

	http.ListenAndServe(":"+port, mux)
}
