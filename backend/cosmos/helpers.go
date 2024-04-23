package cosmos

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

// Helper function for setting up the db connection to a certain endpoint,
// database and container.
func GetContainer() (*azcosmos.ContainerClient, error) {
	connectionString, ok := os.LookupEnv("COSMOS_DB_CONNECTION")
	if !ok {
		log.Fatal("CosmosDbConnection could not be found")
		return nil, errors.New("CosmosDbConnection could not be found")
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		log.Printf("Using default DB_NAME")
		dbName = "gewoscout-db"
	}

	containerName, ok := os.LookupEnv("DB_CONTAINER_NAME")
	if !ok {
		log.Printf("Using default DB_CONTAINER_NAME")
		containerName = "ListingsByCity"
	}

	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		log.Fatal("Failed to create client")
		return nil, err
	}

	container, err := client.NewContainer(dbName, containerName)
	if err != nil {
		log.Fatal("Failed to get container Listings")
		return nil, err
	}
	return container, nil
}

func GetQueryItemsPager(container *azcosmos.ContainerClient, city string, query *models.Query) *runtime.Pager[azcosmos.QueryItemsResponse] {
	var sb strings.Builder
	sb.WriteString("SELECT * FROM c WHERE c._partitionKey = @city")

	queryParams := []azcosmos.QueryParameter{
		{Name: "@city", Value: city},
	}

	if query.MinSize != nil && query.MaxSize != nil {
		queryParams = append(queryParams, azcosmos.QueryParameter{Name: "@minSize", Value: *query.MinSize})
		queryParams = append(queryParams, azcosmos.QueryParameter{Name: "@maxSize", Value: *query.MaxSize})
		sb.WriteString(" AND (c.squareMeters BETWEEN @minSize AND @maxSize)")
	}

	partitionKey := azcosmos.NewPartitionKeyString(strings.ToLower(city))

	options := azcosmos.QueryOptions{
		QueryParameters:   queryParams,
		ContinuationToken: &query.ContinuationToken,
	}

	if query.PageSize != nil {
		options.PageSizeHint = int32(*query.PageSize)
	} else {
		options.PageSizeHint = 10
	}

	return container.NewQueryItemsPager(sb.String(), partitionKey, &options)
}
