import urllib3
from azure.cosmos import CosmosClient, PartitionKey
from azure.storage.queue import QueueServiceClient

from constants import COSMOS_CONNECTION_STRING, QUEUE_CONNECTION_STRING, SCRAPER_RESULTS_QUEUE_NAME, \
    NEW_LISTINGS_QUEUE_NAME

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

if __name__ == "__main__":
    cosmos_client = CosmosClient.from_connection_string(COSMOS_CONNECTION_STRING)

    database_name = 'gewoscout-db'
    container_name = 'ListingsByCity'
    partition_key = PartitionKey(path="/_partitionKey")

    database = cosmos_client.create_database_if_not_exists(id=database_name)
    container = database.create_container_if_not_exists(
        id=container_name,
        partition_key=partition_key
    )

    # Assure the queues are created
    queue_client = QueueServiceClient.from_connection_string(QUEUE_CONNECTION_STRING)
    existing_queues = [x["name"] for x in queue_client.list_queues()]

    expected_queues = [SCRAPER_RESULTS_QUEUE_NAME, NEW_LISTINGS_QUEUE_NAME]
    for queue in expected_queues:
        if queue not in existing_queues:
            queue_client.create_queue(queue)
            print(f"Created queue: {queue}")
