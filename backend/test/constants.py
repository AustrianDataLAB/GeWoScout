import json
import os

CONNECTION_STRING = os.getenv("CONNECTION_STRING", default="AccountEndpoint=https://localhost:8081/;AccountKey=C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw==;")

API_BASE_URL = os.getenv("API_BASE_URL", default="http://localhost:8000/api/")

DIR_PATH = os.path.dirname(os.path.realpath(__file__))

with open(os.path.join(DIR_PATH, 'listings_fixture.json'), 'r') as file:
    LISTINGS_FIXTURE = json.load(file)
