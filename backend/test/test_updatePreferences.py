import pytest
import base64
import json
import requests

from .constants import API_BASE_URL


@pytest.mark.usefixtures(
    "cosmos_db_setup_userdata", "cosmos_db_setup_notification_settings"
)
@pytest.mark.parametrize(
    "city_pref_combo",
    [
        ("vienna", {"city": "vienna", "email": "test@test.com", "minRoomCount": 3}),
        ("graz", {"city": "graz", "email": "test@test.com", "maxRoomCount": 2}),
    ],
)
def test_update_preferences(
    city_pref_combo, cosmos_db_setup_userdata, cosmos_db_setup_notification_settings
):

    principal_data = {
        "auth_typ": "aad",
        "claims": [
            {"typ": "name", "val": "test"},
            {
                "typ": "preferred_username",
                "val": "test@test.com",
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

    user_id = "123123"
    headers = {
        "X-MS-CLIENT-PRINCIPAL-NAME": "shifter",
        "X-MS-CLIENT-PRINCIPAL-ID": user_id,
        "X-MS-CLIENT-PRINCIPAL-IDP": "aad",
        "X-MS-CLIENT-PRINCIPAL": principal_data_base64,
    }

    city, preferences = city_pref_combo
    endpoint_url = f"{API_BASE_URL}/users/preferences"

    response = requests.put(endpoint_url, headers=headers, json=preferences)

    assert response.status_code == 200, "API did not return a successful status code"

    # Check that the entry is also in CosmosDB
    _, container_userdata = cosmos_db_setup_userdata
    item = container_userdata.read_item(item=city, partition_key=user_id)
    assert item is not None, "User data not found in CosmosDB"

    _, container_notification_settings = cosmos_db_setup_notification_settings
    item = container_notification_settings.read_item(item=user_id, partition_key=city)
    assert item is not None, "Notification settings not found in CosmosDB"

    preferences_updated = preferences
    preferences_updated['minSalePrice'] = 300
    response_update = requests.put(endpoint_url, headers=headers, json=preferences_updated)

    assert response.status_code == 200, "API did not return a successful status code"

    # Check that the entry was updated in both 
    _, container_userdata = cosmos_db_setup_userdata
    item = container_userdata.read_item(item=city, partition_key=user_id)
    assert item is not None, "User data not found in CosmosDB"
    assert item['minSalePrice'] == 300

    _, container_notification_settings = cosmos_db_setup_notification_settings
    item = container_notification_settings.read_item(item=user_id, partition_key=city)
    assert item is not None, "Notification settings not found in CosmosDB"
    assert item['minSalePrice'] == 300
