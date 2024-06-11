import base64
import copy
import json
import time

import pytest
from azure.storage.queue import QueueServiceClient

from test.constants import SCRAPER_RESULTS_QUEUE_NAME, NEW_LISTINGS_QUEUE_NAME

TEST_LISTING_MSG_1 = {
    "scraperId": "LinzHousingScraper002",
    "timestamp": "2024-04-06T15:30:00Z",
    "listings": [
        {
            "title": "Modern 3-Bedroom Apartment in Central Linz",
            "housingCooperative": "FutureLiving Genossenschaft",
            "projectId": "FLG2020",
            "listingId": "12345ABC",
            "country": "Austria",
            "city": "Linz",
            "postalCode": "4020",
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
            "previewImageUrl": "https://www.futurelivinggenossenschaft.at/listings/12345ABC/preview.jpg"
        }
    ]
}

TEST_LISTING_1_1 = {
    "id": "futurelivinggenossenschaft_flg2020_12345abc",
    "_partitionKey": "linz",
    "title": "Modern 3-Bedroom Apartment in Central Linz",
    "housingCooperative": "FutureLiving Genossenschaft",
    "projectId": "FLG2020",
    "listingId": "12345ABC",
    "country": "Austria",
    "city": "Linz",
    "postalCode": "4020",
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
    "scraperId": "LinzHousingScraper002",
    "createdAt": "2024-04-06T15:30:00Z",
    "lastModifiedAt": "2024-04-06T15:30:00Z",
}

TEST_LISTING_MSG_2 = {
    "scraperId": "LinzHousingScraper002",
    "timestamp": "2024-04-06T15:30:00Z",
    "listings": [
        {
            "title": "Modern 3-Bedroom Apartment in Central Linz",
            "housingCooperative": "FutureLiving Genossenschaft",
            "projectId": "FLG2020",
            "listingId": "12345ABC",
            "country": "Austria",
            "city": "Linz",
            "postalCode": "4020",
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
            "previewImageUrl": "https://www.futurelivinggenossenschaft.at/listings/12345ABC/preview.jpg"
        },
        {
            "title": "Modern 2-Bedroom Apartment in Central Linz",
            "housingCooperative": "FutureLiving Genossenschaft",
            "projectId": "FLG2020",
            "listingId": "67890DEF",
            "country": "Austria",
            "city": "Ansfelden",
            "postalCode": "4020",
            "address": "Beispielgasse 42",
            "roomCount": 2,
            "squareMeters": 65,
            "availabilityDate": "2024-09-01",
            "yearBuilt": 2019,
            "hwgEnergyClass": "A",
            "fgeeEnergyClass": "A+",
            "listingType": "both",
            "rentPricePerMonth": 800,
            "cooperativeShare": 5000,
            "salePrice": 250000,
            "additionalFees": 6500,
            "detailsUrl": "https://www.futurelivinggenossenschaft.at/listings/67890DEF",
            "previewImageUrl": "https://www.futurelivinggenossenschaft.at/listings/67890DEF/preview.jpg"
        }
    ]
}

TEST_LISTING_2_1 = {
    "id": "futurelivinggenossenschaft_flg2020_12345abc",
    "_partitionKey": "linz",
    "title": "Modern 3-Bedroom Apartment in Central Linz",
    "housingCooperative": "FutureLiving Genossenschaft",
    "projectId": "FLG2020",
    "listingId": "12345ABC",
    "country": "Austria",
    "city": "Linz",
    "postalCode": "4020",
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
    "scraperId": "LinzHousingScraper002",
    "createdAt": "2024-04-06T15:30:00Z",
    "lastModifiedAt": "2024-04-06T15:30:00Z",
}

TEST_LISTING_2_2 = {
    "id": "futurelivinggenossenschaft_flg2020_67890def",
    "_partitionKey": "ansfelden",
    "title": "Modern 2-Bedroom Apartment in Central Linz",
    "housingCooperative": "FutureLiving Genossenschaft",
    "projectId": "FLG2020",
    "listingId": "67890DEF",
    "country": "Austria",
    "city": "Ansfelden",
    "postalCode": "4020",
    "address": "Beispielgasse 42",
    "roomCount": 2,
    "squareMeters": 65,
    "availabilityDate": "2024-09-01",
    "yearBuilt": 2019,
    "hwgEnergyClass": "A",
    "fgeeEnergyClass": "A+",
    "listingType": "both",
    "rentPricePerMonth": 800,
    "cooperativeShare": 5000,
    "salePrice": 250000,
    "additionalFees": 6500,
    "detailsUrl": "https://www.futurelivinggenossenschaft.at/listings/67890DEF",
    "previewImageUrl": "https://www.futurelivinggenossenschaft.at/listings/67890DEF/preview.jpg",
    "scraperId": "LinzHousingScraper002",
    "createdAt": "2024-04-06T15:30:00Z",
    "lastModifiedAt": "2024-04-06T15:30:00Z",
}


@pytest.mark.usefixtures("clear_queues", "queue_service_client", "cosmos_db_setup_listings")
def test_new_listing(queue_service_client: QueueServiceClient, cosmos_db_setup_listings):
    # Add a new listing to the queue
    queue_client = queue_service_client.get_queue_client(SCRAPER_RESULTS_QUEUE_NAME)
    msg = base64.b64encode(json.dumps(TEST_LISTING_MSG_1).encode()).decode()
    queue_client.send_message(msg)

    # Wait for the message to be processed - there should be a new listing in the new listings queue
    new_listings_queue_client = queue_service_client.get_queue_client(NEW_LISTINGS_QUEUE_NAME)
    msgs = read_messages_with_timeout(new_listings_queue_client, 10, 1)
    assert len(msgs) == 1, "Expected 1 message in the new listings queue"

    new_listing = json.loads(base64.b64decode(msgs[0]).decode())
    assert new_listing == TEST_LISTING_1_1, "Expected listing does not match the new listing"

    # Check that the entry is also in CosmosDB
    _, container = cosmos_db_setup_listings
    item = container.read_item(item=TEST_LISTING_1_1['id'], partition_key="linz")
    assert item is not None, "Listing not found in CosmosDB"

    for key, value in TEST_LISTING_1_1.items():
        assert key in item, f"Key {key} not found in listing"
        assert item[key] == value, f"Value for key {key} does not match"


@pytest.mark.usefixtures("clear_queues", "queue_service_client", "cosmos_db_setup_listings")
def test_update_listing(queue_service_client: QueueServiceClient, cosmos_db_setup_listings):
    _, container = cosmos_db_setup_listings
    msg = copy.deepcopy(TEST_LISTING_1_1)
    msg['createdAt'] = "2024-04-06T15:00:00Z"
    msg['lastModifiedAt'] = "2024-04-06T15:00:00Z"
    # TODO update list of patched values
    msg['availabilityDate'] = "ASDF"
    msg['listingType'] = "ASDF"
    msg['rentPricePerMonth'] = 1234
    msg['cooperativeShare'] = 1234
    msg['salePrice'] = 1234
    msg['additionalFees'] = 1234
    msg['detailsUrl'] = "ASDF"

    container.upsert_item(msg)

    # Add an updated listing to the queue
    queue_client = queue_service_client.get_queue_client(SCRAPER_RESULTS_QUEUE_NAME)
    msg = base64.b64encode(json.dumps(TEST_LISTING_MSG_1).encode()).decode()
    queue_client.send_message(msg)

    # Wait for the message to be processed - there should be a new listing in the new listings queue
    new_listings_queue_client = queue_service_client.get_queue_client(NEW_LISTINGS_QUEUE_NAME)
    msgs = read_messages_with_timeout(new_listings_queue_client, 10, 1)
    assert len(msgs) == 0, "Expected 0 messages in the new listings queue"

    # Check that the entry is also in CosmosDB
    _, container = cosmos_db_setup_listings
    item = container.read_item(item=TEST_LISTING_1_1['id'], partition_key="linz")
    assert item is not None, "Listing not found in CosmosDB"

    for key, value in TEST_LISTING_1_1.items():
        if key == 'createdAt':
            assert item[key] == "2024-04-06T15:00:00Z", "Value for key createdAt does not match"
            continue
        assert key in item, f"Key {key} not found in listing"
        assert item[key] == value, f"Value for key {key} does not match"


@pytest.mark.usefixtures("clear_queues", "queue_service_client", "cosmos_db_setup_listings")
def test_multiple_listings(queue_service_client: QueueServiceClient, cosmos_db_setup_listings):
    # One new, one updated
    _, container = cosmos_db_setup_listings

    msg = copy.deepcopy(TEST_LISTING_2_1)
    msg['createdAt'] = "2024-04-06T15:00:00Z"
    msg['lastModifiedAt'] = "2024-04-06T15:00:00Z"
    # TODO update list of patched values
    msg['availabilityDate'] = "ASDF"
    msg['listingType'] = "ASDF"
    msg['rentPricePerMonth'] = 1234
    msg['cooperativeShare'] = 1234
    msg['salePrice'] = 1234
    msg['additionalFees'] = 1234
    msg['detailsUrl'] = "ASDF"

    container.upsert_item(msg)

    # Add an updated and a new listing to the queue
    queue_client = queue_service_client.get_queue_client(SCRAPER_RESULTS_QUEUE_NAME)
    msg = base64.b64encode(json.dumps(TEST_LISTING_MSG_2).encode()).decode()
    queue_client.send_message(msg)

    # Wait for the message to be processed - there should be a new listing in the new listings queue
    new_listings_queue_client = queue_service_client.get_queue_client(NEW_LISTINGS_QUEUE_NAME)
    msgs = read_messages_with_timeout(new_listings_queue_client, 10, 1)
    assert len(msgs) == 1, "Expected 1 message in the new listings queue"

    new_listing = json.loads(base64.b64decode(msgs[0]).decode())
    assert new_listing == TEST_LISTING_2_2, "Expected listing does not match the new listing"

    # Check that the updated entry is also in CosmosDB
    _, container = cosmos_db_setup_listings
    item = container.read_item(item=TEST_LISTING_2_1['id'], partition_key=TEST_LISTING_2_1['_partitionKey'])
    assert item is not None, "Listing 2_1 not found in CosmosDB"

    for key, value in TEST_LISTING_2_1.items():
        if key == 'createdAt':
            assert item[key] == "2024-04-06T15:00:00Z", "Value for key createdAt does not match"
            continue
        assert key in item, f"Key {key} not found in listing"
        assert item[key] == value, f"Value for key {key} does not match"

    # Check that the new entry is also in CosmosDB
    item = container.read_item(item=TEST_LISTING_2_2['id'], partition_key=TEST_LISTING_2_2['_partitionKey'])
    assert item is not None, "Listing 2_2 not found in CosmosDB"

    for key, value in TEST_LISTING_2_2.items():
        assert key in item, f"Key {key} not found in listing"
        assert item[key] == value, f"Value for key {key} does not match"


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
