import requests

from .constants import API_BASE_URL


def test_health():
    response = requests.get(f"{API_BASE_URL}/health")
    assert response.status_code == 200, f"Expected 200, got {response.status_code}"
    assert response.headers["Content-Type"] == "application/json; charset=utf-8", \
        f"Expected 'application/json; charset=utf-8', got {response.headers['Content-Type']}"
    assert response.json() == {"status": "ok"}, f"Expected {{'status': 'ok'}}, got {response.json()}"
