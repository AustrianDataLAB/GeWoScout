import json

CONNECTION_STRING = "<your_connection_string>"

API_BASE_URL = "http://localhost:8000/api/"

with open('listings_fixture.json', 'r') as file:
    LISTINGS_FIXTURE = json.load(file)
