import logging

import pytest
import urllib3
from azure.cosmos import PartitionKey, CosmosClient, exceptions
from azure.cosmos.documents import ConnectionPolicy

from .constants import LISTINGS_FIXTURE

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

logging.basicConfig(level=logging.DEBUG)


@pytest.fixture(scope="module")
def cosmos_db_client():
    # client = CosmosClient.from_connection_string(CONNECTION_STRING)
    connection_policy = ConnectionPolicy()
    connection_policy.RequestTimeout = 10
    connection_policy.ConnectionMode = 1
    connection_policy.DisableSSLVerification = True
    connection_policy.EnableEndpointDiscovery = True
    # connection_policy.PreferredLocations = ["https://localhost:8081"]

    client = CosmosClient(
        url="https://localhost:8081/",
        credential="C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw==",
        connection_policy=connection_policy,
    )

    print(client.list_databases())
    # client.client_connection.ReadEndpoint = client.client_connection.url_connection
    # client.client_connection.WriteEndpoint = client.client_connection.url_connection
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
    # cosmos_db_client.delete_database(database_name)

    # Cleanup: remove all items from the container
    for item in container.read_all_items():
        container.delete_item(item['id'], partition_key=item['_partitionKey'])
