package cosmos

import (
	"errors"
	"log"
	"os"

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
