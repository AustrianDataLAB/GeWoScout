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
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type ReturnValue struct {
	Data string
}

type InvokeResponse struct {
	Outputs     map[string]interface{}
	Logs        []string
	ReturnValue interface{}
}

type InvokeRequest struct {
	Data     map[string]interface{}
	Metadata map[string]interface{}
}

// Represents a listing as it is queried from various Genossenschaft pages.
type Listing struct {
	Id string `json:"id"`
}

// Holds any form of error (either from Azure or some internal error)
type Error struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func queueTriggerHandler(w http.ResponseWriter, r *http.Request) {
	var invokeReq InvokeRequest

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&invokeReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: implement something
	md, _ := json.Marshal(invokeReq.Metadata)
	d, _ := json.Marshal(invokeReq.Data)
	log.Printf("Received metadata: %s\n", md)
	log.Printf("Received data: %s\n", d)

	invokeResponse := InvokeResponse{Logs: []string{}, ReturnValue: nil}

	resp, err := json.Marshal(invokeResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

type CosmosData struct {
	Documents string `json:"documents"`
}

type CosmosMetadataSys struct {
	MethodName string `json:"MethodName"`
	UtcNow     string `json:"UtcNow"`
	RandGuid   string `json:"RandGuid"`
}

type CosmosMetadata struct {
	Sys CosmosMetadataSys `json:"sys"`
}

type CosmosTrigger struct {
	Data     CosmosData     `json:"Data"`
	Metadata CosmosMetadata `json:"Metadata"`
}

func cosmosUpdateHandler(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	fmt.Println(t.Month())
	fmt.Println(t.Day())
	fmt.Println(t.Year())
	fmt.Println(r.Header)
	ua := r.Header.Get("User-Agent")
	fmt.Printf("user agent is: %s \n", ua)
	invocationid := r.Header.Get("X-Azure-Functions-InvocationId")
	fmt.Printf("invocationid is: %s \n", invocationid)

	queryParams := r.URL.Query()

	for k, v := range queryParams {
		fmt.Println("k:", k, "v:", v)
	}

	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	fmt.Println("Body:", string(body))

	var cosmosTrigger CosmosTrigger
	err := json.Unmarshal(body, &cosmosTrigger)
	if err != nil {
		fmt.Println("Error unmarshalling JSON trigger: ", err)
	}

	fmt.Println("CosmosTrigger:", cosmosTrigger)

	fmt.Println("Cosmos Trigger:")
	fmt.Println("Metadata", cosmosTrigger.Metadata.Sys.MethodName, cosmosTrigger.Metadata.Sys.UtcNow, cosmosTrigger.Metadata.Sys.RandGuid)
	fmt.Println("Data", cosmosTrigger.Data.Documents)

	var dataStr string
	err = json.Unmarshal([]byte(cosmosTrigger.Data.Documents), &dataStr)
	if err != nil {
		fmt.Println("Error unmarshalling JSON string: ", err)
	}

	// `MyCosmosDocument` contains an escaped JSON string, parse it too.
	var documents []map[string]interface{} // Using a map for flexible structure handling
	err = json.Unmarshal([]byte(dataStr), &documents)
	if err != nil {
		fmt.Println("Error unmarshalling JSON documents: ", err)
	}

	fmt.Printf("Parsed Documents: %+v\n", documents)

	returnValue := ReturnValue{Data: documents[0]["id"].(string)}
	invokeResponse := InvokeResponse{Logs: []string{"test log1", "test log2"}, ReturnValue: returnValue}

	js, err := json.Marshal(invokeResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Helper function for setting up the db connection to a certain endpoint,
// database and container.
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

// Handler function for /listings, which returns any listings within the
// partition defined by the city path param, which is guaranteed to exist at
// this point.
func getListings(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")

	container, err := getContainer()
	if err != nil {
		render.JSON(w, r, Error{Message: err.Error(), StatusCode: http.StatusInternalServerError})
		return
	}

	query := "SELECT * FROM c WHERE c._partitionKey = @city"
	partitionKey := azcosmos.NewPartitionKeyString(strings.ToLower(city))
	options := azcosmos.QueryOptions{QueryParameters: []azcosmos.QueryParameter{{Name: "@city", Value: city}}}
	pager := container.NewQueryItemsPager(query, partitionKey, &options)

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

	render.JSON(w, r, listings)
}

// Handler function for /listings/{listingId}, which fetches a specific
// listing within the partition defined by the city key.
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

func SetupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	r.Get("/QueueTrigger", queueTriggerHandler)
	r.Get("/CosmosTrigger", cosmosUpdateHandler)
	r.Get("/api/listings/{city}", getListings)
	r.Get("/api/listings/{city}/{listingId}", getListingsById)

	return r
}

func main() {
	port, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !exists {
		port = "8080"
	}

	log.Printf("About to listen on %s. Go to http://127.0.0.1:%s/", port, port)

	r := SetupRouter()

	http.ListenAndServe(":"+port, r)
}
