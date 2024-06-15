import pytest
import base64
import json
import requests

from .constants import API_BASE_URL, USERDATA_PREFERENCES_FIXTURE


@pytest.mark.usefixtures(
    "cosmos_db_setup_userdata", "cosmos_db_setup_notification_settings"
)
@pytest.mark.parametrize(
    "userdata_prefs",
    USERDATA_PREFERENCES_FIXTURE,
)
def test_get_preferences(
    userdata_prefs, cosmos_db_setup_userdata, cosmos_db_setup_notification_settings
):
    email = userdata_prefs['preferences'][0]['email']

    principal_data = {
        "auth_typ": "aad",
        "claims": [
            {"typ": "name", "val": "test"},
            {
                "typ": "preferred_username",
                "val": email,
            },
            {"typ": "uti", "val": "odhfaosdhfodsah"},
            {"typ": "ver", "val": "2.0"},
        ],
        "name_typ": "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name",
        "role_typ": "http://schemas.microsoft.com/ws/2008/06/identity/claims/role",
    }

    principal_data_base64 = base64.b64encode(
        json.dumps(principal_data).encode()
    ).decode()

    user_id = userdata_prefs['user_id'] 
    headers = {
        "X-MS-CLIENT-PRINCIPAL-NAME": email,
        "X-MS-CLIENT-PRINCIPAL-ID": user_id,
        "X-MS-CLIENT-PRINCIPAL-IDP": "aad",
        "X-MS-CLIENT-PRINCIPAL": principal_data_base64,
    }

    endpoint_url = f"{API_BASE_URL}/users/preferences"

    response = requests.get(endpoint_url, headers=headers)

    assert response.status_code == 200, "API did not return a successful status code"
    assert (
        response.headers["Content-Type"] == "application/json; charset=utf-8"
    ), "API did not return JSON"

    response_json = response.json()

    assert len(response_json) == len(userdata_prefs['preferences'])
    assert response_json == userdata_prefs['preferences']
