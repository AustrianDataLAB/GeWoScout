import os
import subprocess
import time
import warnings

import pytest
import requests
import urllib3
from azure.core.exceptions import ResourceNotFoundError
from azure.cosmos import PartitionKey, CosmosClient, exceptions
from azure.storage.queue import QueueServiceClient

from .constants import (
    COSMOS_CONNECTION_STRING,
    LISTINGS_FIXTURE,
    USERDATA_PREFERENCES_FIXTURE,
    QUEUE_CONNECTION_STRING,
    LISTING_PREFERENCE_MATCH_INPUT_QUEUE_NAME,
    LISTING_PREFERENCE_MATCH_OUTPUT_QUEUE_NAME,
    QUEUE_NAMES,
)


@pytest.fixture(scope="session", autouse=True)
def disable_warnings():
    urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
    warnings.filterwarnings("ignore", category=DeprecationWarning)


@pytest.fixture(scope="module")
def cosmos_db_client():
    client = CosmosClient.from_connection_string(COSMOS_CONNECTION_STRING)
    return client


@pytest.fixture(scope="module")
def queue_service_client():
    return QueueServiceClient.from_connection_string(QUEUE_CONNECTION_STRING)


@pytest.fixture(scope="module")
def listing_preference_match_input_queue(queue_service_client):
    try:
        queue_client = queue_service_client.get_queue_client(
            LISTING_PREFERENCE_MATCH_INPUT_QUEUE_NAME
        )
        queue_client.clear_messages()
    except ResourceNotFoundError:
        queue_client = queue_service_client.create_queue(
            LISTING_PREFERENCE_MATCH_INPUT_QUEUE_NAME
        )

    yield queue_client

    queue_client.clear_messages()


@pytest.fixture(scope="module")
def listing_preference_match_output_queue(queue_service_client):
    try:
        queue_client = queue_service_client.get_queue_client(
            LISTING_PREFERENCE_MATCH_OUTPUT_QUEUE_NAME
        )
        queue_client.clear_messages()
    except ResourceNotFoundError:
        queue_client = queue_service_client.create_queue(
            LISTING_PREFERENCE_MATCH_OUTPUT_QUEUE_NAME
        )

    yield queue_client

    queue_client.clear_messages()


@pytest.fixture(scope="module")
def clear_queues(queue_service_client):
    queues = QUEUE_NAMES

    for queue_name in queues:
        queue_client = queue_service_client.get_queue_client(queue_name)
        queue_client.clear_messages()

    yield

    for queue_name in queues:
        queue_client = queue_service_client.get_queue_client(queue_name)
        queue_client.clear_messages()


@pytest.fixture(scope="module")
def cosmos_db_setup_listings(cosmos_db_client):
    # Create the container if it does not exist
    container_name = "ListingsByCity"
    database, container = _cosmos_db_setup(cosmos_db_client, container_name)

    for listing in LISTINGS_FIXTURE:
        container.upsert_item(listing)

    # Provide the database and container to the test function
    yield database, container

    # Cleanup: delete database
    # cosmos_db_client.delete_database(database_name)

    # Cleanup: remove all items from the container
    for item in container.read_all_items():
        container.delete_item(item["id"], partition_key=item["_partitionKey"])


@pytest.fixture(scope="module")
def cosmos_db_setup_userdata(cosmos_db_client):
    # Create the container if it does not exist
    container_name = "UserDataByUserId"
    database, container = _cosmos_db_setup(cosmos_db_client, container_name)

    for ud in USERDATA_PREFERENCES_FIXTURE:
        for preference in ud["preferences"]:
            container.upsert_item(preference)

    # Provide the database and container to the test function
    yield database, container

    # Cleanup: delete database
    # cosmos_db_client.delete_database(database)

    # Cleanup: remove all items from the container
    for item in container.read_all_items():
        container.delete_item(item["id"], partition_key=item["_partitionKey"])


@pytest.fixture(scope="module")
def cosmos_db_setup_notification_settings(cosmos_db_client):
    # Create the container if it does not exist
    container_name = "NotificationSettingsByCity"
    database, container = _cosmos_db_setup(cosmos_db_client, container_name)

    for ud in USERDATA_PREFERENCES_FIXTURE:
        for i in range(len(ud["preferences"])):
            pref = ud["preferences"][i].copy()
            pref["_partitionKey"] = pref["city"]
            pref["id"] = ud["user_id"]
            container.upsert_item(pref)

    # Provide the database and container to the test function
    yield database, container

    # Cleanup: delete database
    # cosmos_db_client.delete_database(database_name)

    # Cleanup: remove all items from the container
    for item in container.read_all_items():
        container.delete_item(item["id"], partition_key=item["_partitionKey"])


def _cosmos_db_setup(cosmos_db_client, container_name):
    # Create the database if it does not exist
    database_name = "gewoscout-db"
    try:
        database = cosmos_db_client.create_database_if_not_exists(id=database_name)
    except exceptions.CosmosHttpResponseError:
        database = cosmos_db_client.get_database_client(database_name)

    partition_key = PartitionKey(path="/_partitionKey")
    try:
        container = database.create_container_if_not_exists(
            id=container_name, partition_key=partition_key
        )
    except exceptions.CosmosHttpResponseError:
        container = database.get_container_client(container_name)
    return database, container


@pytest.fixture(scope="session", autouse=True)
def setup_backend_server():
    # Start the server in the background
    print("Starting the server...")

    file_dir = os.path.dirname(os.path.realpath(__file__))
    parent_dir = os.path.dirname(file_dir)
    print("Working directory: ", parent_dir)

    proc = subprocess.Popen(
        ["func", "start", "--port", "8000"], shell=True, cwd=parent_dir
    )

    try:
        # Wait for the server to be up by checking the health endpoint
        timeout = 60
        start_time = time.time()
        url = "http://localhost:8000/api/health"

        while True:
            try:
                response = requests.get(url)
                if response.status_code == 200:
                    print("Server is up and running!")
                    break
            except requests.ConnectionError:
                pass  # The server is not up yet

            if time.time() - start_time > timeout:
                raise TimeoutError("Timed out waiting for the server to start")
            time.sleep(1)  # Wait a bit before trying again

    except Exception as e:
        print(f"An error occurred: {e}")
        proc.kill()  # Ensure server is killed if setup fails
        raise

    yield

    # Teardown: stop the server
    print("Stopping the server...")
    proc.terminate()
    proc.wait()
