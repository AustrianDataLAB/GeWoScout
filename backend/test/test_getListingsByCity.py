import urllib

import pytest
import requests

from .constants import API_BASE_URL, LISTINGS_FIXTURE


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "city,expected_listings_count", [("vienna", 7), ("graz", 1), ("salzburg", 1)]
)
def test_get_all_listings_for_city(city, expected_listings_count):
    endpoint_url = f"{API_BASE_URL}/cities/{city}/listings"
    expected_ids = [
        listing["id"]
        for listing in LISTINGS_FIXTURE
        if listing["_partitionKey"] == city
    ]

    # Send a GET request to the endpoint
    response = requests.get(endpoint_url)

    # Check if the request was successful
    assert response.status_code == 200, "API did not return a successful status code"
    # Check that content type is JSON
    assert (
        response.headers["Content-Type"] == "application/json; charset=utf-8"
    ), "API did not return JSON"

    # Parse JSON response
    response_json = response.json()
    assert isinstance(response_json, dict), "Response format is not a dictionary"

    assert "results" in response_json, "Response does not contain 'results' key"
    listings = response_json["results"]

    # Assuming you know the exact number of listings the fixture should have for Vienna
    assert (
        len(listings) == expected_listings_count
    ), f"Expected {expected_listings_count} listings, got {len(listings)}"
    ids = [listing["id"] for listing in listings]
    assert ids == expected_ids, "Listing IDs do not match"

    # Additional checks could include verifying that all listings are for Vienna
    for listing in listings:
        assert listing.get("_partitionKey") == city, "Listing city mismatch"

    assert (
        "continuationToken" in response_json
    ), "Response does not contain 'continuationToken' key"
    assert (
        response_json["continuationToken"] is None
    ), "Expected continuation token to be None"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
def test_get_listings_using_continuation_token():
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?pageSize=3"

    # Initial page
    response1 = requests.get(endpoint_url)

    assert response1.status_code == 200, "API did not return a successful status code"
    assert (
        response1.headers["Content-Type"] == "application/json; charset=utf-8"
    ), "API did not return JSON"

    response1_json = response1.json()
    assert isinstance(response1_json, dict), "Response format is not a dictionary"

    assert "results" in response1_json, "Response does not contain 'results' key"
    results_1 = response1_json["results"]

    assert len(results_1) == 3, "Expected 3 listings in the first page"
    ids_1 = [listing["id"] for listing in results_1]
    assert ids_1 == [
        "FutureLivingGenossenschaft_FLG2024_12345ABC",
        "UrbanLiving_ULProj001_ULFlat001",
        "UrbanLiving_ULProj002_ULFlat004",
    ], "Wrong listing IDs in the first page"

    assert (
        "continuationToken" in response1_json
    ), "Response does not contain 'continuationToken' key"
    continuation_token = response1_json["continuationToken"]

    # Next page
    response2 = requests.get(
        f"{endpoint_url}&continuationToken={urllib.parse.quote(continuation_token)}"
    )

    assert response2.status_code == 200, "API did not return a successful status code"
    assert (
        response2.headers["Content-Type"] == "application/json; charset=utf-8"
    ), "API did not return JSON"

    response2_json = response2.json()
    assert isinstance(response2_json, dict), "Response format is not a dictionary"

    assert "results" in response2_json, "Response does not contain 'results' key"
    results_2 = response2_json["results"]

    assert len(results_2) == 3, "Expected 3 listings in the second page"
    ids_2 = [listing["id"] for listing in results_2]
    assert ids_2 == [
        "CityHomes_CHProj204_CHFlat059",
        "CityHomes_CHProj205_CHFlat062",
        "GreenLiving_GLProj987_GLFlat321",
    ], "Wrong listing IDs in the second page"

    assert (
        "continuationToken" in response2_json
    ), "Response does not contain 'continuationToken' key"
    continuation_token = response2_json["continuationToken"]

    # Next page
    response3 = requests.get(
        f"{endpoint_url}&continuationToken={urllib.parse.quote(continuation_token)}"
    )

    assert response3.status_code == 200, "API did not return a successful status code"
    assert (
        response3.headers["Content-Type"] == "application/json; charset=utf-8"
    ), "API did not return JSON"

    response3_json = response3.json()
    assert isinstance(response3_json, dict), "Response format is not a dictionary"

    assert "results" in response3_json, "Response does not contain 'results' key"
    results_3 = response3_json["results"]

    assert len(results_3) == 1, "Expected 1 listing in the third page"
    ids_3 = [listing["id"] for listing in results_3]
    assert ids_3 == [
        "ArtHabitat_AHProj053_AHFlat678"
    ], "Wrong listing IDs in the third page"

    assert (
        "continuationToken" in response3_json
    ), "Response does not contain 'continuationToken' key"
    continuation_token = response3_json["continuationToken"]
    assert (
        continuation_token is None
    ), "Expected continuation token to be None in the last page"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
def test_get_listings_for_nonexistent_city():
    endpoint_url = f"{API_BASE_URL}/cities/berlin/listings"

    response = requests.get(endpoint_url)

    assert response.status_code == 200, "API did not return a successful status code"
    assert (
        response.headers["Content-Type"] == "application/json; charset=utf-8"
    ), "API did not return JSON"

    response_json = response.json()
    assert isinstance(response_json, dict), "Response format is not a dictionary"

    assert "results" in response_json, "Response does not contain 'results' key"
    assert (
        len(response_json["results"]) == 0
    ), "Expected 0 listings for a nonexistent city"

    assert (
        "continuationToken" in response_json
    ), "Response does not contain 'continuationToken' key"
    assert (
        response_json["continuationToken"] is None
    ), "Expected continuation token to be None"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
def test_get_listings_has_default_page_size_10(cosmos_db_setup_listings):
    _database, container = cosmos_db_setup_listings

    # Insert 20 listings for Vienna
    for i in range(20):
        container.upsert_item({"id": f"listing{i}", "_partitionKey": "vienna"})

    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings"
    response = requests.get(endpoint_url)

    assert response.status_code == 200, "API did not return a successful status code"
    response_json = response.json()

    assert len(response_json["results"]) == 10, "Expected 10 listings by default"
    assert (
        response_json["continuationToken"] is not None
    ), "Expected continuation token to be present"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
def test_get_listings_query_min_flat_size():
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?minSqm=80"
    expected_ids = [
        l["id"]
        for l in LISTINGS_FIXTURE
        if l["_partitionKey"] == "vienna" and l["squareMeters"] >= 80
    ]

    response = requests.get(endpoint_url)

    assert response.status_code == 200, "API did not return a successful status code"
    response_json = response.json()

    assert len(response_json["results"]) == len(
        expected_ids
    ), "Expected number of listings does not match"
    assert [
        l["id"] for l in response_json["results"]
    ] == expected_ids, "Listing IDs do not match"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
def test_get_listings_query_max_flat_size():
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?maxSqm=80"
    expected_ids = [
        l["id"]
        for l in LISTINGS_FIXTURE
        if l["_partitionKey"] == "vienna" and l["squareMeters"] <= 80
    ]

    response = requests.get(endpoint_url)

    assert response.status_code == 200, "API did not return a successful status code"
    response_json = response.json()

    assert len(response_json["results"]) == len(
        expected_ids
    ), "Expected number of listings does not match"
    assert [
        l["id"] for l in response_json["results"]
    ] == expected_ids, "Listing IDs do not match"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
def test_get_listings_query_min_max_flat_size():
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?minSqm=80.5&maxSqm=100"
    expected_ids = [
        l["id"]
        for l in LISTINGS_FIXTURE
        if l["_partitionKey"] == "vienna" and 80.5 <= l["squareMeters"] <= 100
    ]

    response = requests.get(endpoint_url)

    assert response.status_code == 200, "API did not return a successful status code"
    response_json = response.json()

    assert len(response_json["results"]) == len(
        expected_ids
    ), "Expected number of listings does not match"
    assert [
        l["id"] for l in response_json["results"]
    ] == expected_ids, "Listing IDs do not match"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
def test_get_listings_with_invalid_city_returns_error():
    endpoint_url = f"{API_BASE_URL}/cities/%20/listings"

    response = requests.get(endpoint_url)

    assert response.status_code == 400, "API did not return a 400 status code"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "continuation_token",
    [
        "invalid",
        "%2BRIDLEL%3A~9TZrALpFOV8DAAAAAAAAAA%3D%3D%23RT%3A1%23TRC%3A2%23ISV%3A2%23IEO%3A65567%23QCF%3A8%23FPC%3AAgEAAFFFAAHA%2BAI%3D",
        "+RID:~9TZrALpFOV8DAAAAAAAAAA==#RT:1#TRC:2#ISV:2#IEO:65567#QCF:8#FPC:AgEAAAAEAAHApAI=",
    ],
)
def test_get_listings_invalid_continuation_token_returns_error(continuation_token):
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?pageSize=2&continuationToken={continuation_token}"

    response = requests.get(endpoint_url)

    assert response.status_code == 400, "API did not return a 400 status code"


@pytest.mark.parametrize("page_size", [0, -1, 31])  # Max page size is 30
def test_get_listings_invalid_page_size_returns_error(page_size):
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?pageSize={page_size}"

    response = requests.get(endpoint_url)

    assert response.status_code == 400, "API did not return a 400 status code"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "size_query",
    [
        "minSqm=invalid",
        "minSqm=-1",
        "minSqm=inf",
        "maxSqm=0",
        "maxSqm=invalid",
        "minSqm=100&maxSqm=50",
        "minSqm=51&maxSqm=50",
    ],
)
def test_get_listings_invalid_size_query_returns_error(size_query):
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?{size_query}"

    response = requests.get(endpoint_url)

    assert response.status_code == 400, "API did not return a 400 status code"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "endpoint, expected_ids",
    [
        (
            f"{API_BASE_URL}/cities/vienna/listings?title=3-Bedroom",
            [
                l["id"]
                for l in LISTINGS_FIXTURE
                if l["_partitionKey"] == "vienna" and "3-Bedroom" in l["title"]
            ],
        ),
        (
            f"{API_BASE_URL}/cities/vienna/listings?housingCooperative=Urban%20Living",
            [
                l["id"]
                for l in LISTINGS_FIXTURE
                if l["_partitionKey"] == "vienna"
                and l.get("housingCooperative") == "Urban Living"
            ],
        ),
        (
            f"{API_BASE_URL}/cities/vienna/listings?projectId=FLG2024",
            [
                l["id"]
                for l in LISTINGS_FIXTURE
                if l["_partitionKey"] == "vienna" and l.get("projectId") == "FLG2024"
            ],
        ),
        (
            f"{API_BASE_URL}/cities/vienna/listings?postalCode=1010,1050",
            [
                l["id"]
                for l in LISTINGS_FIXTURE
                if l["_partitionKey"] == "vienna"
                and l.get("postalCode") in ["1010", "1050"]
            ],
        ),
    ],
)
def test_get_listings_basic_query(endpoint, expected_ids):
    response = requests.get(endpoint)

    assert response.status_code == 200, "API did not return a successful status code"
    response_json = response.json()

    assert len(response_json["results"]) == len(
        expected_ids
    ), "Expected number of listings does not match"
    assert [
        l["id"] for l in response_json["results"]
    ] == expected_ids, "Listing IDs do not match"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "endpoint, min_room_count, max_room_count, room_count, expected_ids",
    [
        (
            f"{API_BASE_URL}/cities/vienna/listings?roomCount=3",
            None,
            None,
            3,
            [
                l["id"]
                for l in LISTINGS_FIXTURE
                if l["_partitionKey"] == "vienna" and l.get("roomCount") == 3
            ],
        ),
        (
            f"{API_BASE_URL}/cities/vienna/listings?minRoomCount=2",
            2,
            None,
            None,
            [
                l["id"]
                for l in LISTINGS_FIXTURE
                if l["_partitionKey"] == "vienna" and l.get("roomCount", 0) >= 2
            ],
        ),
        (
            f"{API_BASE_URL}/cities/vienna/listings?maxRoomCount=4",
            None,
            4,
            None,
            [
                l["id"]
                for l in LISTINGS_FIXTURE
                if l["_partitionKey"] == "vienna" and l.get("roomCount", 0) <= 4
            ],
        ),
        (
            f"{API_BASE_URL}/cities/vienna/listings?minRoomCount=2&maxRoomCount=4",
            2,
            4,
            None,
            [
                l["id"]
                for l in LISTINGS_FIXTURE
                if l["_partitionKey"] == "vienna" and 2 <= l.get("roomCount", 0) <= 4
            ],
        ),
    ],
)
def test_get_listings_room_count_query(
    endpoint, min_room_count, max_room_count, room_count, expected_ids
):
    response = requests.get(endpoint)

    assert response.status_code == 200, "API did not return a successful status code"
    response_json = response.json()

    assert len(response_json["results"]) == len(
        expected_ids
    ), "Expected number of listings does not match"
    assert [
        l["id"] for l in response_json["results"]
    ] == expected_ids, "Listing IDs do not match"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "min_room_count, max_room_count, expected_status",
    [
        (-1, None, 400),
        (None, -1, 400),
        (-1, -1, 400),
        (3, 2, 400),
        (2, 3, 200),
        (2, None, 200),
        (None, 2, 200),
    ],
)
def test_get_listings_room_count_validation(
    min_room_count, max_room_count, expected_status
):
    query_params = []
    if min_room_count is not None:
        query_params.append(f"minRoomCount={min_room_count}")
    if max_room_count is not None:
        query_params.append(f"maxRoomCount={max_room_count}")
    query_string = "&".join(query_params)
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?{query_string}"

    response = requests.get(endpoint_url)

    assert (
        response.status_code == expected_status
    ), f"API did not return expected status code {expected_status}"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "min_sqm, max_sqm, expected_status",
    [
        (100, 80, 400),
        (80, 100, 200),
    ],
)
def test_get_listings_sqm_validation(min_sqm, max_sqm, expected_status):
    query_string = f"minSqm={min_sqm}&maxSqm={max_sqm}"
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?{query_string}"

    response = requests.get(endpoint_url)

    assert (
        response.status_code == expected_status
    ), f"API did not return expected status code {expected_status}"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "min_year_built, max_year_built, expected_status",
    [
        (2010, 2000, 400),
        (2000, 2010, 200),
    ],
)
def test_get_listings_year_built_validation(
    min_year_built, max_year_built, expected_status
):
    query_string = f"minYearBuilt={min_year_built}&maxYearBuilt={max_year_built}"
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?{query_string}"

    response = requests.get(endpoint_url)

    assert (
        response.status_code == expected_status
    ), f"API did not return expected status code {expected_status}"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "date, expected_status",
    [
        ("2024-01-01", 200),
        ("invalid-date", 400),
    ],
)
def test_get_listings_available_from_validation(date, expected_status):
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?availableFrom={date}"
    response = requests.get(endpoint_url)
    assert (
        response.status_code == expected_status
    ), f"API did not return expected status code {expected_status}"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "min_class, expected_status",
    [
        ("A++", 200),
        ("A+", 200),
        ("A", 200),
        ("B", 200),
        ("C", 200),
        ("D", 200),
        ("E", 200),
        ("F", 200),
        ("invalid-class", 400),
    ],
)
def test_get_listings_hwg_energy_class_validation(min_class, expected_status):
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings"
    params = {
        "minHwgEnergyClass": min_class,
    }
    response = requests.get(endpoint_url, params=params)
    assert (
        response.status_code == expected_status
    ), f"API did not return expected status code {expected_status}"
    if expected_status == 200:
        response_json = response.json()
        valid_classes = ["A++", "A+", "A", "B", "C", "D", "E", "F", "G", None]
        max_index = valid_classes.index(min_class)
        assert all(
            (
                l["hwgEnergyClass"] in valid_classes[: max_index + 1]
                if min_class is not None and min_class != "G"
                else valid_classes + [""]
            )
            for l in response_json["results"]
        ), f"Listing energy classes do not match the required minimum, should match {matching_classes}"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "min_class, expected_status",
    [
        ("A++", 200),
        ("A+", 200),
        ("A", 200),
        ("B", 200),
        ("C", 200),
        ("D", 200),
        ("E", 200),
        ("F", 200),
        ("invalid-class", 400),
    ],
)
def test_get_listings_fgee_energy_class_validation(min_class, expected_status):
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings"
    params = {
        "minFgeeEnergyClass": min_class,
    }
    response = requests.get(endpoint_url, params=params)
    assert (
        response.status_code == expected_status
    ), f"API did not return expected status code {expected_status}"
    if expected_status == 200:
        response_json = response.json()
        valid_classes = ["A++", "A+", "A", "B", "C", "D", "E", "F", "G", None]
        max_index = valid_classes.index(min_class)
        assert all(
            l["fgeeEnergyClass"]
            in (
                valid_classes[: max_index + 1]
                if min_class is not None and min_class != "G"
                else valid_classes + [""]
            )
            for l in response_json["results"]
        ), "Listing energy classes do not match the required minimum"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "listing_type, expected_status",
    [
        ("rent", 200),
        ("sale", 200),
        ("both", 200),
        ("invalid-type", 400),
    ],
)
def test_get_listings_listing_type_validation(listing_type, expected_status):
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?listingType={listing_type}"
    response = requests.get(endpoint_url)
    assert (
        response.status_code == expected_status
    ), f"API did not return expected status code {expected_status}"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "param, value, expected_status",
    [
        ("minRentPrice", -1, 400),
        ("minRentPrice", 100, 200),
        ("maxRentPrice", -1, 400),
        ("maxRentPrice", 1000, 200),
    ],
)
def test_get_listings_rent_price_validation(param, value, expected_status):
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?{param}={value}"
    response = requests.get(endpoint_url)
    assert (
        response.status_code == expected_status
    ), f"API did not return expected status code {expected_status}"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "min_rent, max_rent, expected_status",
    [
        (1000, 500, 400),
        (500, 1000, 200),
    ],
)
def test_get_listings_rent_price_range_validation(min_rent, max_rent, expected_status):
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?minRentPrice={min_rent}&maxRentPrice={max_rent}"
    response = requests.get(endpoint_url)
    assert (
        response.status_code == expected_status
    ), f"API did not return expected status code {expected_status}"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "param, value, expected_status",
    [
        ("minCooperativeShare", -1, 400),
        ("minCooperativeShare", 100, 200),
        ("maxCooperativeShare", -1, 400),
        ("maxCooperativeShare", 1000, 200),
    ],
)
def test_get_listings_cooperative_share_validation(param, value, expected_status):
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?{param}={value}"
    response = requests.get(endpoint_url)
    assert (
        response.status_code == expected_status
    ), f"API did not return expected status code {expected_status}"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "min_share, max_share, expected_status",
    [
        (1000, 500, 400),
        (500, 1000, 200),
    ],
)
def test_get_listings_cooperative_share_range_validation(
    min_share, max_share, expected_status
):
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?minCooperativeShare={min_share}&maxCooperativeShare={max_share}"
    response = requests.get(endpoint_url)
    assert (
        response.status_code == expected_status
    ), f"API did not return expected status code {expected_status}"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "param, value, expected_status",
    [
        ("minSalePrice", -1, 400),
        ("minSalePrice", 100000, 200),
        ("maxSalePrice", -1, 400),
        ("maxSalePrice", 500000, 200),
    ],
)
def test_get_listings_sale_price_validation(param, value, expected_status):
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?{param}={value}"
    response = requests.get(endpoint_url)
    assert (
        response.status_code == expected_status
    ), f"API did not return expected status code {expected_status}"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "min_price, max_price, expected_status",
    [
        (500000, 100000, 400),
        (100000, 500000, 200),
    ],
)
def test_get_listings_sale_price_range_validation(
    min_price, max_price, expected_status
):
    endpoint_url = f"{API_BASE_URL}/cities/vienna/listings?minSalePrice={min_price}&maxSalePrice={max_price}"
    response = requests.get(endpoint_url)
    assert (
        response.status_code == expected_status
    ), f"API did not return expected status code {expected_status}"


@pytest.mark.usefixtures("cosmos_db_setup_listings")
@pytest.mark.parametrize(
    "sort_by, sort_type, expected_sorted_field, expected_status",
    [
        ("squareMeters", "ASC", "squareMeters", 200),
        ("squareMeters", "DESC", "squareMeters", 200),
        ("roomCount", "ASC", "roomCount", 200),
        ("roomCount", "DESC", "roomCount", 200),
        ("yearBuilt", "ASC", "yearBuilt", 200),
        ("yearBuilt", "DESC", "yearBuilt", 200),
        ("roomCount", "LEL", None, 400),
    ],
)
def test_get_listings_sort_by(
    sort_by, sort_type, expected_sorted_field, expected_status
):
    endpoint_url = (
        f"{API_BASE_URL}/cities/vienna/listings?sortBy={sort_by}&sortType={sort_type}"
    )
    response = requests.get(endpoint_url)

    assert (
        response.status_code == expected_status
    ), f"API did not return expected status code {expected_status}"

    if expected_status == 200:
        response_json = response.json()
        listings = response_json["results"]
        if sort_type == "ASC":
            assert listings == sorted(
                listings, key=lambda x: x.get(expected_sorted_field is None, float("inf"))
            ), f"Listings are not sorted in ascending order by {sort_by}"
        elif sort_type == "DESC":
            assert listings == sorted(
                listings,
                key=lambda x: x.get(expected_sorted_field is None, float("-inf")),
                reverse=True,
            ), f"Listings are not sorted in descending order by {sort_by}"
