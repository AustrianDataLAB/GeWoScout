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

	container, err := client.NewContainer(dbName, "ListingsByCity")
	if err != nil {
		log.Fatal("Failed to get container Listings")
		return nil, err
	}
	return container, nil
}
