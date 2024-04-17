import json

import pytest
from azure.cosmos import CosmosClient, PartitionKey, exceptions

from secrets import CONNECTION_STRING


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

    # Load and insert data
    with open('listings_fixture.json', 'r') as file:
        listings_data = json.load(file)

    for listing in listings_data:
        container.upsert_item(listing)

    # Provide the database and container to the test function
    yield database, container

    # Cleanup: delete database
    cosmos_db_client.delete_database(database_name)


@pytest.mark.usefixtures("cosmos_db_setup")
def test_cosmosdb_fixture_data_insertion(cosmos_db_setup):
    _database, container = cosmos_db_setup

    # Check if data insertion was correct
    # There are 7 listings in the fixture for Vienna
    query = "SELECT * FROM c WHERE c._partitionKey = 'vienna'"
    items = list(container.query_items(query=query, enable_cross_partition_query=True))
    assert len(items) == 7

    # Should also work this way
    items = list(container.query_items(query=query, partition_key='vienna'))
    assert len(items) == 7

    query = "SELECT * FROM c WHERE c._partitionKey = 'salzburg'"
    items = list(container.query_items(query=query, partition_key='salzburg'))
    assert len(items) == 1

# Additional tests can be defined similarly
