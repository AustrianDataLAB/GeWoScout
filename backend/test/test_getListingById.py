import pytest
import requests

from .constants import API_BASE_URL, LISTINGS_FIXTURE


@pytest.mark.usefixtures("cosmos_db_setup")
@pytest.mark.parametrize("id_comb", [
    ("vienna", "ArtHabitat_AHProj053_AHFlat678"),
    ("salzburg", "ArtHabitat_AHProj055_AHFlat681"),
    ("graz", "GreenLiving_GLProj990_GLFlat325")
])
def test_get_listing(id_comb):
    city, listing_id = id_comb
    endpoint_url = f"{API_BASE_URL}cities/{city}/listings/{listing_id}"
    expected_listing = [l for l in LISTINGS_FIXTURE if l["id"] == listing_id][0]

    response = requests.get(endpoint_url)

    assert response.status_code == 200, "API did not return a successful status code"
    assert response.headers["Content-Type"] == "application/json; charset=utf-8", "API did not return JSON"

    response_data = response.json()
    assert response_data["id"] == listing_id
    assert response_data["_partitionKey"] == city

    for key in expected_listing:
        assert response_data[key] == expected_listing[key]


@pytest.mark.usefixtures("cosmos_db_setup")
@pytest.mark.parametrize("id_comb", [
    ("unknown", "ArtHabitat_AHProj055_AHFlat681"),  # unknown city
    ("salzburg", "unknown"),  # unknown listing
    ("vienna", "GreenLiving_GLProj990_GLFlat325"),  # pk mismatch
    ("graz", "ArtHabitat_AHProj055_AHFlat681")  # pk mismatch
])
def test_get_listing_not_found(id_comb):
    city, listing_id = id_comb
    endpoint_url = f"{API_BASE_URL}cities/{city}/listings/{listing_id}"

    response = requests.get(endpoint_url)

    assert response.status_code == 404, "API did not return a 404 status code"
