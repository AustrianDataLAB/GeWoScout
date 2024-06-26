import json
import os

COSMOS_CONNECTION_STRING = os.getenv("CONNECTION_STRING", default="AccountEndpoint=https://localhost:8081/;AccountKey=C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw==;")
QUEUE_CONNECTION_STRING = os.getenv("QUEUE_CONNECTION_STRING", default="DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;QueueEndpoint=http://127.0.0.1:10001/devstoreaccount1;")

API_BASE_URL = os.getenv("API_BASE_URL", default="http://localhost:8000/api")

DIR_PATH = os.path.dirname(os.path.realpath(__file__))

SCRAPER_RESULTS_QUEUE_NAME = "scraper-result-queue"
NEW_LISTINGS_QUEUE_NAME = "new-listings-queue"
LISTING_PREFERENCE_MATCH_INPUT_QUEUE_NAME = "listing-preference-match-input-queue"
LISTING_PREFERENCE_MATCH_OUTPUT_QUEUE_NAME = "listing-preference-match-output-queue"

QUEUE_NAMES = [
    SCRAPER_RESULTS_QUEUE_NAME,
    NEW_LISTINGS_QUEUE_NAME,
    LISTING_PREFERENCE_MATCH_INPUT_QUEUE_NAME,
    LISTING_PREFERENCE_MATCH_OUTPUT_QUEUE_NAME,
]

with open(os.path.join(DIR_PATH, 'listings_fixture.json'), 'r') as file:
    LISTINGS_FIXTURE = json.load(file)

with open(os.path.join(DIR_PATH, 'userdata_preferences_fixture.json'), 'r') as file:
    USERDATA_PREFERENCES_FIXTURE = json.load(file)

