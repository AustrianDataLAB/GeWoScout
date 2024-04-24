import requests

from .constants import API_BASE_URL


def test_base_redirect():
    endpoint_url = f"{API_BASE_URL}/swagger"

    response = requests.get(endpoint_url, allow_redirects=False)

    assert response.status_code == 301, "API did not return a 301 status code"
    assert response.headers[
               "Location"] == "/api/swagger/index.html", "API did not redirect to the correct URL"


def test_base_redirect_slash():
    endpoint_url = f"{API_BASE_URL}/swagger/"

    response = requests.get(endpoint_url, allow_redirects=False)

    assert response.status_code == 301, "API did not return a 301 status code"
    assert response.headers[
               "Location"] == "/api/swagger/index.html", "API did not redirect to the correct URL"


def test_swagger_main_html():
    endpoint_url = f"{API_BASE_URL}/swagger/index.html"

    response = requests.get(endpoint_url)

    assert response.status_code == 200, "API did not return a successful status code"
    assert response.headers["Content-Type"] == "text/html; charset=utf-8", "API did not return HTML"

    assert response.text.lstrip().startswith("<!DOCTYPE html>"), "API did not return the expected HTML content"


def test_swagger_doc_json():
    endpoint_url = f"{API_BASE_URL}/swagger/doc.json"

    response = requests.get(endpoint_url)

    assert response.status_code == 200, "API did not return a successful status code"
    assert response.headers["Content-Type"] == "application/json; charset=utf-8", "API did not return JSON"

    response_data = response.json()
    assert response_data["info"]["title"] == "GeWoScout API"
