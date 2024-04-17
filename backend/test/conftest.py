import pytest
from azure.cosmos import PartitionKey, CosmosClient, exceptions

from .constants import CONNECTION_STRING, LISTINGS_FIXTURE


@pytest.fixture(scope="module")
def cosmos_db_client():
    client = CosmosClient.from_connection_string(CONNECTION_STRING)
    return client


@pytest.fixture(scope="module")
def cosmos_db_setup(cosmos_db_client):
    # Create the database if it does not exist
    database_name = 'gewoscout-db'
    try:
        database = cosmos_db_client.create_database_if_not_exists(id=database_name)
    except exceptions.CosmosHttpResponseError:
        database = cosmos_db_client.get_database_client(database_name)

    # Create the container if it does not exist
    container_name = 'ListingsByCity'
    partition_key = PartitionKey(path="/_partitionKey")
    try:
        container = database.create_container_if_not_exists(
            id=container_name,
            partition_key=partition_key
        )
    except exceptions.CosmosHttpResponseError:
        container = database.get_container_client(container_name)

    for listing in LISTINGS_FIXTURE:
        container.upsert_item(listing)

    # Provide the database and container to the test function
    yield database, container

    # Cleanup: delete database
    cosmos_db_client.delete_database(database_name)
