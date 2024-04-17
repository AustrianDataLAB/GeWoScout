import json
import os

CONNECTION_STRING = "<your_connection_string>"

API_BASE_URL = "http://localhost:8000/api/"

DIR_PATH = os.path.dirname(os.path.realpath(__file__))

with open(os.path.join(DIR_PATH, 'listings_fixture.json'), 'r') as file:
    LISTINGS_FIXTURE = json.load(file)
