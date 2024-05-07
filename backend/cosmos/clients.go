package cosmos

import (
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"log"
	"os"
)

func InitClients() (client *azcosmos.Client, db *azcosmos.DatabaseClient, container *azcosmos.ContainerClient) {
	connectionString, ok := os.LookupEnv("COSMOS_DB_CONNECTION")
	if !ok {
		log.Fatal("CosmosDbConnection could not be found")
	}

	var err error
	client, err = azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		log.Fatal("Failed to create client")
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		log.Printf("Using default DB_NAME")
		dbName = "gewoscout-db"
	}

	db, err = client.NewDatabase(dbName)
	if err != nil {
		log.Fatal("Failed to get database")
	}

	containerName, ok := os.LookupEnv("DB_CONTAINER_NAME")
	if !ok {
		log.Printf("Using default DB_CONTAINER_NAME")
		containerName = "ListingsByCity"
	}

	container, err = db.NewContainer(containerName)
	if err != nil {
		log.Fatal("Failed to get container Listings")
	}

	return
}
