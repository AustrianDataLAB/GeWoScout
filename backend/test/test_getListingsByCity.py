import urllib

import pytest
import requests

from .constants import API_BASE_URL, LISTINGS_FIXTURE


@pytest.mark.usefixtures("cosmos_db_setup")
@pytest.mark.parametrize("city,expected_listings_count", [
    ("vienna", 7),
    ("graz", 1),
    ("salzburg", 1)
])
def test_get_all_listings_for_city(city, expected_listings_count):
    endpoint_url = f"{API_BASE_URL}/cities/{city}/listings"
    expected_ids = [listing["id"] for listing in LISTINGS_FIXTURE if listing['_partitionKey'] == city]

    # Send a GET request to the endpoint
    response = requests.get(endpoint_url)

    # Check if the request was successful
    assert response.status_code == 200, "API did not return a successful status code"
    # Check that content type is JSON
    assert response.headers['Content-Type'] == 'application/json; charset=utf-8', "API did not return JSON"

    # Parse JSON response
    response_json = response.json()
    assert isinstance(response_json, dict), "Response format is not a dictionary"

    assert "results" in response_json, "Response does not contain 'results' key"
    listings = response_json["results"]

    # Assuming you know the exact number of listings the fixture should have for Vienna
    assert len(listings) == expected_listings_count, f"Expected {expected_listings_count} listings, got {len(listings)}"
    ids = [listing["id"] for listing in listings]
    assert ids == expected_ids, "Listing IDs do not match"

    # Additional checks could include verifying that all listings are for Vienna
    for listing in listings:
        assert listing.get('_partitionKey') == city, "Listing city mismatch"

    assert "continuationToken" in response_json, "Response does not contain 'continuationToken' key"
    assert response_json["continuationToken"] is None, "Expected continuation token to be None"


@pytest.mark.usefixtures("cosmos_db_setup")
def test_get_listings_using_continuation_token():
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?pageSize=3"

    # Initial page
    response1 = requests.get(endpoint_url)

    assert response1.status_code == 200, "API did not return a successful status code"
    assert response1.headers['Content-Type'] == 'application/json; charset=utf-8', "API did not return JSON"

    response1_json = response1.json()
    assert isinstance(response1_json, dict), "Response format is not a dictionary"

    assert "results" in response1_json, "Response does not contain 'results' key"
    results_1 = response1_json["results"]

    assert len(results_1) == 3, "Expected 3 listings in the first page"
    ids_1 = [listing["id"] for listing in results_1]
    assert ids_1 == [
        "FutureLivingGenossenschaft_FLG2024_12345ABC",
        "UrbanLiving_ULProj001_ULFlat001",
        "UrbanLiving_ULProj002_ULFlat004"
    ], "Wrong listing IDs in the first page"

    assert "continuationToken" in response1_json, "Response does not contain 'continuationToken' key"
    continuation_token = response1_json["continuationToken"]

    # Next page
    response2 = requests.get(f"{endpoint_url}&continuationToken={urllib.parse.quote(continuation_token)}")

    assert response2.status_code == 200, "API did not return a successful status code"
    assert response2.headers['Content-Type'] == 'application/json; charset=utf-8', "API did not return JSON"

    response2_json = response2.json()
    assert isinstance(response2_json, dict), "Response format is not a dictionary"

    assert "results" in response2_json, "Response does not contain 'results' key"
    results_2 = response2_json["results"]

    assert len(results_2) == 3, "Expected 3 listings in the second page"
    ids_2 = [listing["id"] for listing in results_2]
    assert ids_2 == [
        "CityHomes_CHProj204_CHFlat059",
        "CityHomes_CHProj205_CHFlat062",
        "GreenLiving_GLProj987_GLFlat321"
    ], "Wrong listing IDs in the second page"

    assert "continuationToken" in response2_json, "Response does not contain 'continuationToken' key"
    continuation_token = response2_json["continuationToken"]

    # Next page
    response3 = requests.get(f"{endpoint_url}&continuationToken={urllib.parse.quote(continuation_token)}")

    assert response3.status_code == 200, "API did not return a successful status code"
    assert response3.headers['Content-Type'] == 'application/json; charset=utf-8', "API did not return JSON"

    response3_json = response3.json()
    assert isinstance(response3_json, dict), "Response format is not a dictionary"

    assert "results" in response3_json, "Response does not contain 'results' key"
    results_3 = response3_json["results"]

    assert len(results_3) == 1, "Expected 1 listing in the third page"
    ids_3 = [listing["id"] for listing in results_3]
    assert ids_3 == [
        "ArtHabitat_AHProj053_AHFlat678"
    ], "Wrong listing IDs in the third page"

    assert "continuationToken" in response3_json, "Response does not contain 'continuationToken' key"
    continuation_token = response3_json["continuationToken"]
    assert continuation_token is None, "Expected continuation token to be None in the last page"


@pytest.mark.usefixtures("cosmos_db_setup")
def test_get_listings_for_nonexistent_city():
    endpoint_url = f"{API_BASE_URL}/cities/berlin/listings"

    response = requests.get(endpoint_url)

    assert response.status_code == 200, "API did not return a successful status code"
    assert response.headers['Content-Type'] == 'application/json; charset=utf-8', "API did not return JSON"

    response_json = response.json()
    assert isinstance(response_json, dict), "Response format is not a dictionary"

    assert "results" in response_json, "Response does not contain 'results' key"
    assert len(response_json["results"]) == 0, "Expected 0 listings for a nonexistent city"

    assert "continuationToken" in response_json, "Response does not contain 'continuationToken' key"
    assert response_json["continuationToken"] is None, "Expected continuation token to be None"


@pytest.mark.usefixtures("cosmos_db_setup")
def test_get_listings_has_default_page_size_10(cosmos_db_setup):
    _database, container = cosmos_db_setup

    # Insert 20 listings for Vienna
    for i in range(20):
        container.upsert_item({
            "id": f"listing{i}",
            "_partitionKey": "vienna"
        })

    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings"
    response = requests.get(endpoint_url)

    assert response.status_code == 200, "API did not return a successful status code"
    response_json = response.json()

    assert len(response_json["results"]) == 10, "Expected 10 listings by default"
    assert response_json["continuationToken"] is not None, "Expected continuation token to be present"


@pytest.mark.usefixtures("cosmos_db_setup")
def test_get_listings_query_min_flat_size():
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?minSize=80"
    expected_ids = [l["id"] for l in LISTINGS_FIXTURE if l['_partitionKey'] == 'vienna' and l['squareMeters'] >= 80]

    response = requests.get(endpoint_url)

    assert response.status_code == 200, "API did not return a successful status code"
    response_json = response.json()

    assert len(response_json["results"]) == len(expected_ids), "Expected number of listings does not match"
    assert [l["id"] for l in response_json["results"]] == expected_ids, "Listing IDs do not match"


@pytest.mark.usefixtures("cosmos_db_setup")
def test_get_listings_query_max_flat_size():
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?maxSize=80"
    expected_ids = [l["id"] for l in LISTINGS_FIXTURE if l['_partitionKey'] == 'vienna' and l['squareMeters'] <= 80]

    response = requests.get(endpoint_url)

    assert response.status_code == 200, "API did not return a successful status code"
    response_json = response.json()

    assert len(response_json["results"]) == len(expected_ids), "Expected number of listings does not match"
    assert [l["id"] for l in response_json["results"]] == expected_ids, "Listing IDs do not match"


@pytest.mark.usefixtures("cosmos_db_setup")
def test_get_listings_query_min_max_flat_size():
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?minSize=80&maxSize=100"
    expected_ids = [l["id"] for l in LISTINGS_FIXTURE if
                    l['_partitionKey'] == 'vienna' and 80 <= l['squareMeters'] <= 100]

    response = requests.get(endpoint_url)

    assert response.status_code == 200, "API did not return a successful status code"
    response_json = response.json()

    assert len(response_json["results"]) == len(expected_ids), "Expected number of listings does not match"
    assert [l["id"] for l in response_json["results"]] == expected_ids, "Listing IDs do not match"


@pytest.mark.usefixtures("cosmos_db_setup")
def test_get_listings_with_invalid_city_returns_error():
    endpoint_url = f"{API_BASE_URL}/cities/%20/listings"

    response = requests.get(endpoint_url)

    assert response.status_code == 400, "API did not return a 400 status code"


@pytest.mark.usefixtures("cosmos_db_setup")
@pytest.mark.parametrize("continuation_token", [
    "invalid",
    "%2BRIDLEL%3A~9TZrALpFOV8DAAAAAAAAAA%3D%3D%23RT%3A1%23TRC%3A2%23ISV%3A2%23IEO%3A65567%23QCF%3A8%23FPC%3AAgEAAFFFAAHA%2BAI%3D",
    "+RID:~9TZrALpFOV8DAAAAAAAAAA==#RT:1#TRC:2#ISV:2#IEO:65567#QCF:8#FPC:AgEAAAAEAAHA+AI="
])
def test_get_listings_invalid_continuation_token_returns_error(continuation_token):
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?pageSize=2&continuationToken={continuation_token}"

    response = requests.get(endpoint_url)

    assert response.status_code == 400, "API did not return a 400 status code"


@pytest.mark.parametrize("page_size", [
    0,
    -1,
    31  # Max page size is 30
])
def test_get_listings_invalid_page_size_returns_error(page_size):
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?pageSize={page_size}"

    response = requests.get(endpoint_url)

    assert response.status_code == 400, "API did not return a 400 status code"


@pytest.mark.usefixtures("cosmos_db_setup")
@pytest.mark.parametrize("size_query", [
    "minSize=invalid",
    "minSize=-1",
    "minSize=inf",
    "maxSize=0",
    "maxSize=invalid",
    "minSize=100&maxSize=50",
    "minSize=51&maxSize=50",
])
def test_get_listings_invalid_size_query_returns_error(size_query):
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?{size_query}"

    response = requests.get(endpoint_url)

    assert response.status_code == 400, "API did not return a 400 status code"
