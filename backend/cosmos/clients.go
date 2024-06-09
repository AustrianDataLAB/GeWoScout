package cosmos

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func InitClients() (client *azcosmos.Client, db *azcosmos.DatabaseClient) {
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

	return
}

func InitContainerClient(db *azcosmos.DatabaseClient, containerName string) *azcosmos.ContainerClient {
	container, err := db.NewContainer(containerName)
	if err != nil {
		log.Fatal("Failed to get container Listings")
	}
	return container
}
