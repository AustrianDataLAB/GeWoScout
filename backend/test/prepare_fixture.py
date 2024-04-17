from azure.cosmos import CosmosClient, PartitionKey

from constants import CONNECTION_STRING

if __name__ == "__main__":
    client = CosmosClient.from_connection_string(CONNECTION_STRING)

    database_name = 'gewoscout-db'
    container_name = 'ListingsByCity'
    partition_key = PartitionKey(path="/_partitionKey")

    database = client.create_database_if_not_exists(id=database_name)
    container = database.create_container_if_not_exists(
        id=container_name,
        partition_key=partition_key
    )
