import json
import urllib.parse

import pytest
import requests
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


API_BASE_URL = "http://localhost:8000/api/"


@pytest.mark.usefixtures("cosmos_db_setup")
@pytest.mark.parametrize("city,expected_listings_count", [
    ("vienna", 7),
    ("graz", 1),
    ("salzburg", 1)
])
def test_get_all_listings_for_city(city, expected_listings_count):
    endpoint_url = f"{API_BASE_URL}cities/{city}/listings"

    # Send a GET request to the endpoint
    response = requests.get(endpoint_url)

    # Check if the request was successful
    assert response.status_code == 200, "API did not return a successful status code"
    # Check that content type is JSON
    assert response.headers['Content-Type'] == 'application/json', "API did not return JSON"

    # Parse JSON response
    listings = response.json()
    assert isinstance(listings, list), "Response format is not a list"

    # Assuming you know the exact number of listings the fixture should have for Vienna
    assert len(listings) == expected_listings_count, f"Expected {expected_listings_count} listings, got {len(listings)}"

    # Additional checks could include verifying that all listings are for Vienna
    for listing in listings:
        assert listing.get('_partitionKey') == city, "Listing city mismatch"


@pytest.mark.usefixtures("cosmos_db_setup")
def test_get_listings_using_continuation_token():
    endpoint_url = f"{API_BASE_URL}cities/vienna/listings?pageSize=3"

    # Initial page
    response1 = requests.get(endpoint_url)

    assert response1.status_code == 200, "API did not return a successful status code"
    assert response1.headers['Content-Type'] == 'application/json', "API did not return JSON"

    response1_json = response1.json()
    assert isinstance(response1_json, dict), "Response format is not a dictionary"

    assert "results" in response1_json, "Response does not contain 'results' key"
    results_1 = response1_json["results"]

    assert len(results_1) == 3, "Expected 3 listings in the first page"

    assert "continuationToken" in response1_json, "Response does not contain 'continuationToken' key"
    continuation_token = response1_json["continuationToken"]

    # Next page
    response2 = requests.get(f"{endpoint_url}&continuationToken={urllib.parse.quote(continuation_token)}")

    assert response2.status_code == 200, "API did not return a successful status code"
    assert response2.headers['Content-Type'] == 'application/json', "API did not return JSON"

    response2_json = response2.json()
    assert isinstance(response2_json, dict), "Response format is not a dictionary"

    assert "results" in response2_json, "Response does not contain 'results' key"
    results_2 = response2_json["results"]

    assert len(results_2) == 3, "Expected 3 listings in the second page"

    assert "continuationToken" in response2_json, "Response does not contain 'continuationToken' key"
    continuation_token = response2_json["continuationToken"]

    # Next page
    response3 = requests.get(f"{endpoint_url}&continuationToken={urllib.parse.quote(continuation_token)}")

    assert response3.status_code == 200, "API did not return a successful status code"
    assert response3.headers['Content-Type'] == 'application/json', "API did not return JSON"

    response3_json = response3.json()
    assert isinstance(response3_json, dict), "Response format is not a dictionary"

    assert "results" in response3_json, "Response does not contain 'results' key"
    results_3 = response3_json["results"]

    assert len(results_3) == 1, "Expected 1 listing in the third page"

    assert "continuationToken" in response3_json, "Response does not contain 'continuationToken' key"
    continuation_token = response3_json["continuationToken"]
    assert continuation_token is None, "Expected continuation token to be None in the last page"
