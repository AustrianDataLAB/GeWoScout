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
    "title": "3-Bedroom",
    "housingCooperative": "FutureLiving",
    "projectId": "FLG2024",
    "postalCode": "1010",
    "minRoomCount": 2,
    "maxRoomCount": 4,
    "minSqm": 50,
    "maxSqm": 100,
    "availableFrom": "2024-01-01",
    "minYearBuilt": 1990,
    "maxYearBuilt": 2020,
    "minHwgEnergyClass": "A",
    "minFgeeEnergyClass": "A+",
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

# MinSqm is 96, but the listing has 95
PREF_3 = {
    "id": "1233",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "minSqm": 96,
    "email": "test3@test.com",
}

# MaxSqm is 94, but the listing has 95
PREF_4 = {
    "id": "1234",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "maxSqm": 94.5,
    "email": "test4@test.com",
}

# MinRoomCount is 4, but the listing has 3
PREF_5 = {
    "id": "1235",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "minRoomCount": 4,
    "email": "test5@test.com"
}

# MaxRoomCount is 2, but the listing has 3
PREF_6 = {
    "id": "1236",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "maxRoomCount": 2,
    "email": "test6@test.com"
}

# Title is not in the listing
PREF_7 = {
    "id": "1237",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "title": "NOT_IN_TITLE",
    "email": "test7@test.com"
}

# HousingCooperative is not in the listing
PREF_8 = {
    "id": "1238",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "housingCooperative": "NOT_IN_HOUSING_COOP",
    "email": "test8@test.com"
}

# ProjectId wrong
PREF_9 = {
    "id": "1239",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "projectId": "WRONG_PROJECT_ID",
    "email": "test9@test.com"
}

# PostalCode wrong
PREF_10 = {
    "id": "1240",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "postalCode": "WRONG_POSTAL_CODE",
    "email": "test10@test.com"
}

# AvailabilityFrom wrong
PREF_11 = {
    "id": "1241",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "availableFrom": "2024-10-01",
    "email": "test11@test.com"
}

# MinYearBuilt too high
PREF_12 = {
    "id": "1242",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "minYearBuilt": 2020,
    "email": "test12@test.com"
}

# MaxYearBuilt too low
PREF_13 = {
    "id": "1243",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "maxYearBuilt": 2018,
    "email": "test13@test.com"
}

# MinHwgEnergyClass too high
PREF_14 = {
    "id": "1244",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "minHwgEnergyClass": "A+",
    "email": "test14@test.com"
}

# MinFgeeEnergyClass too low
PREF_15 = {
    "id": "1245",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "minFgeeEnergyClass": "A++",
    "email": "test15@test.com"
}

# MinRentPrice too high
PREF_16 = {
    "id": "1246",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "minRentPrice": 1201,
    "email": "test16@test.com"
}

# MaxRentPrice too low
PREF_17 = {
    "id": "1247",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "maxRentPrice": 1199,
    "email": "test17@test.com"
}

# MinCooperativeShare too high
PREF_18 = {
    "id": "1248",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "minCooperativeShare": 5001,
    "email": "test18@test.com"
}

# MaxCooperativeShare too low
PREF_19 = {
    "id": "1249",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "maxCooperativeShare": 4999,
    "email": "test19@test.com"
}

# MinSalePrice too high
PREF_20 = {
    "id": "1250",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "minSalePrice": 350001,
    "email": "test20@test.com"
}

# MaxSalePrice too low
PREF_21 = {
    "id": "1251",
    "_partitionKey": "vienna",
    "city": "Vienna",
    "maxSalePrice": 349999,
    "email": "test21@test.com"
}

# TODO add more tests for listingType

PREFS = [
    PREF_1,
    PREF_2,
    PREF_3,
    PREF_4,
    PREF_5,
    PREF_6,
    PREF_7,
    PREF_8,
    PREF_9,
    PREF_10,
    PREF_11,
    PREF_12,
    PREF_13,
    PREF_14,
    PREF_15,
    PREF_16,
    PREF_17,
    PREF_18,
    PREF_19,
    PREF_20,
    PREF_21,
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
    "listing_preference_match_input_queue",
    "listing_preference_match_output_queue",
)
def test_get_preferences(listing_preference_match_input_queue, listing_preference_match_output_queue, cosmos_db_setup_notification_settings):
    _, container = cosmos_db_setup_notification_settings

    for pref in PREFS:
        container.upsert_item(pref)

    # Add a new listing to the queue
    msg = base64.b64encode(json.dumps(LISTING_1).encode()).decode()
    listing_preference_match_input_queue.send_message(msg)

    # Wait for the message to be processed - there should be a new listing in the new listings queue
    msgs = read_messages_with_timeout(listing_preference_match_output_queue, 10, 1)
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
