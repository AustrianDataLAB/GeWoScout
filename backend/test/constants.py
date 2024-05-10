import json
import os

COSMOS_CONNECTION_STRING = os.getenv("CONNECTION_STRING", default="AccountEndpoint=https://localhost:8081/;AccountKey=C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw==;")
QUEUE_STORAGE_CONNECTION = os.getenv("QUEUE_STORAGE_CONNECTION", default="DefaultEndpointsProtocol=https;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;QueueEndpoint=https://127.0.0.1:10001/devstoreaccount1;")

API_BASE_URL = os.getenv("API_BASE_URL", default="http://localhost:8000/api")

DIR_PATH = os.path.dirname(os.path.realpath(__file__))

SCRAPER_RESULTS_QUEUE_NAME = "scraper-result-queue"
NEW_LISTINGS_QUEUE_NAME = "new-listings-queue"

with open(os.path.join(DIR_PATH, 'listings_fixture.json'), 'r') as file:
    LISTINGS_FIXTURE = json.load(file)
