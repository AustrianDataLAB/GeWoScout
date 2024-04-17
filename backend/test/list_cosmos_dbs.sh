#!/bin/bash

# Cosmos DB account settings
account="your_cosmos_db_account"
key="C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="

# Prepare the headers and verb for the request
verb="get"
resourcetype="colls"
resourceid="dbs/gewoscout-db"
date=$(TZ=GMT date "+%a, %d %b %Y %H:%M:%S %Z")
version="2018-12-31"

# Generate the authorization header
string_to_sign=$(printf "%s\n%s\n%s\n%s\n\n" "$verb" "$resourcetype" "$resourceid" "$date" | tr -d '\r')
signature=$(echo -n "$string_to_sign" | openssl dgst -sha256 -hmac $(echo -n "$key" | base64 -d) -binary | base64 | tr '+/' '-_')

auth_header="type=master&ver=1.0&sig=$signature"

# Execute the curl command to list all containers in the database "gewoscout-db"
curl -X GET "https://localhost:8081/dbs/gewoscout-db/colls" \
    -H "Authorization: $auth_header" \
    -H "x-ms-date: $date" \
    -H "x-ms-version: $version" \
    -H "Content-Type: application/json"
i
