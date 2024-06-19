import pytest
import base64
import json
import requests

from azure.cosmos.exceptions import CosmosResourceNotFoundError

from .constants import API_BASE_URL, USERDATA_PREFERENCES_FIXTURE


@pytest.mark.usefixtures(
    "cosmos_db_setup_userdata", "cosmos_db_setup_notification_settings"
)
@pytest.mark.parametrize(
    "userdata_prefs",
    USERDATA_PREFERENCES_FIXTURE,
)
def test_delete_preferences(
    userdata_prefs, cosmos_db_setup_userdata, cosmos_db_setup_notification_settings
):
    _, container_userdata = cosmos_db_setup_userdata
    _, container_notification_settings = cosmos_db_setup_notification_settings

    preferences = userdata_prefs["preferences"]

    for preference in preferences:
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

        user_id = userdata_prefs["user_id"]
        headers = {
            "X-MS-CLIENT-PRINCIPAL-NAME": "tester",
            "X-MS-CLIENT-PRINCIPAL-ID": user_id,
            "X-MS-CLIENT-PRINCIPAL-IDP": "aad",
            "X-MS-CLIENT-PRINCIPAL": principal_data_base64,
        }

        city = preference["city"]
        endpoint_url = f"{API_BASE_URL}/users/preferences/{city}"

        ud_item_orig = container_userdata.read_item(item=city, partition_key=user_id)
        assert ud_item_orig is not None, "User data not found in CosmosDB"
        ns_item_orig = container_notification_settings.read_item(
            item=user_id, partition_key=city
        )
        assert ns_item_orig is not None, "Notification settings not found in CosmosDB"

        response = requests.delete(endpoint_url, headers=headers)

        assert (
            response.status_code == 200
        ), "API did not return a successful status code"

        # Check that the entry is also removed in CosmosDB
        with pytest.raises(CosmosResourceNotFoundError):
            ud_item = container_userdata.read_item(item=city, partition_key=user_id)

        with pytest.raises(CosmosResourceNotFoundError):
            # Check that the entry is also removed in CosmosDB
            ns_item = container_notification_settings.read_item(
                item=user_id, partition_key=city
            )
