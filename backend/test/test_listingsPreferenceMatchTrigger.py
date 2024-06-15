import base64
import json
import time

import pytest
from azure.storage.queue import QueueServiceClient

from .constants import LISTING_PREFERENCE_MATCH_INPUT_QUEUE_NAME, LISTING_PREFERENCE_MATCH_OUTPUT_QUEUE_NAME

PREF_1 = {
    "id": "1231",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "email": "test1@test.com",
    "title": "Test Title",
    "housingCooperative": "Test Cooperative",
    "projectId": "12345",
    "postalCode": "123456",
    "roomCount": 3,
    "minRoomCount": 2,
    "maxRoomCount": 4,
    "minSqm": 50,
    "maxSqm": 100,
    "availableFrom": "2024-01-01",
    "minYearBuilt": 1990,
    "maxYearBuilt": 2020,
    "minHwgEnergyClass": "A",
    "minFgeeEnergyClass": "B",
    "listingType": "rent",
    "minRentPrice": 500,
    "maxRentPrice": 1500,
    "minCooperativeShare": 1000,
    "maxCooperativeShare": 10000,
    "minSalePrice": 200000
}

PREF_2 = {
    "id": "1232",
    "_partitionKey": "vienna",
    "city": "vienna",
    "minSqm": 50,
    "email": "test2@test.com",
}

PREF_3 = {
    "id": "1233",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "minSqm": 250,
    "email": "test3@test.com",
}

PREF_4 = {
    "id": "1234",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "maxSqm": 60,
    "email": "test4@test.com",
}

PREFS = [
    PREF_1,
    PREF_2,
    PREF_3,
]

LISTING_1 = {
    "id": "FutureLivingGenossenschaft_FLG2024_12345ABC",
    "_partitionKey": "vienna",
    "title": "Modern 3-Bedroom Apartment in Central Vienna",
    "housingCooperative": "FutureLiving Genossenschaft",
    "projectId": "FLG2024",
    "listingId": "12345ABC",
    "country": "Austria",
    "city": "Vienna",
    "postalCode": "1010",
    "address": "Beispielgasse 42",
    "roomCount": 3,
    "squareMeters": 95,
    "availabilityDate": "2024-09-01",
    "yearBuilt": 2019,
    "hwgEnergyClass": "A",
    "fgeeEnergyClass": "A+",
    "listingType": "both",
    "rentPricePerMonth": 1200,
    "cooperativeShare": 5000,
    "salePrice": 350000,
    "additionalFees": 6500,
    "detailsUrl": "https://www.futurelivinggenossenschaft.at/listings/12345ABC",
    "previewImageUrl": "https://www.futurelivinggenossenschaft.at/listings/12345ABC/preview.jpg",
    "scraperId": "viennaHousingScraper002",
    "createdAt": "2024-04-06T15:30:00Z",
    "lastModifiedAt": "2024-04-06T15:30:00Z"
}


@pytest.mark.usefixtures(
    # "cosmos_db_setup_notification_settings",
    "clear_queues",
    "queue_service_client"
)
def test_get_preferences(queue_service_client: QueueServiceClient, cosmos_db_setup_notification_settings):
    print("HI")

    _, container = cosmos_db_setup_notification_settings

    for pref in PREFS:
        container.upsert_item(pref)

    # Add a new listing to the queue
    queue_client = queue_service_client.get_queue_client(LISTING_PREFERENCE_MATCH_INPUT_QUEUE_NAME)
    msg = base64.b64encode(json.dumps(LISTING_1).encode()).decode()
    queue_client.send_message(msg)

    # Wait for the message to be processed - there should be a new listing in the new listings queue
    output_queue_client = queue_service_client.get_queue_client(LISTING_PREFERENCE_MATCH_OUTPUT_QUEUE_NAME)
    msgs = read_messages_with_timeout(output_queue_client, 10, 1)
    assert len(msgs) == 1, "Expected 1 message in the new listings queue"

    new_notification_data = json.loads(base64.b64decode(msgs[0]).decode())
    assert new_notification_data["listing"] == LISTING_1, "Expected the new listing to be the same as the one sent"

    expected_mails = sorted([
        "test1@test.com",
        "test2@test.com"
    ])
    assert sorted(new_notification_data["emails"]) == expected_mails, \
        "Expected the emails to be the same as the ones in the preferences"


def read_messages_with_timeout(queue_client, timeout, max_messages):
    start_time = time.time()
    messages = []

    while True:
        received_messages = queue_client.receive_messages(max_messages=max_messages)
        messages.extend(x.content for x in received_messages)

        # Check if timeout has been reached or we have received the desired number of messages
        elapsed_time = time.time() - start_time
        if elapsed_time > timeout or len(messages) >= max_messages:
            break

    return messages
