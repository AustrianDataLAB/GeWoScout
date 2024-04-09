import base64
import binascii
import json
import logging
import os

from azure.cosmos import CosmosClient, exceptions
import azure.functions as func


def main(req: func.HttpRequest) -> func.HttpResponse:

    try:
        # First is Azure Functions, second is local development
        cosmosdb_connection_string = os.environ["CUSTOMCONNSTR_CosmosDBConnectionString"] if "CUSTOMCONNSTR_CosmosDBConnectionString" in os.environ \
            else os.environ["ConnectionStrings:CosmosDBConnectionString"]
        logging.error(f"Connection string: {cosmosdb_connection_string}")

        client = CosmosClient.from_connection_string(cosmosdb_connection_string)

        database_name = 'gewoscout-db'
        container_name = 'ListingsByCity'

        partition_key = req.route_params.get('city').lower().replace(" ", "") if 'city' in req.route_params else ''
        if partition_key == "":
            return func.HttpResponse(
                "City parameter is required",
                status_code=400
            )

        page_size_param = req.params.get('pageSize') if 'pageSize' in req.params else "10"

        try:
            page_size = int(page_size_param)
        except ValueError:
            return func.HttpResponse(
                "Invalid pageSize parameter",
                status_code=400
            )
        
        if page_size < 1:
            return func.HttpResponse(
                "pageSize must be greater than 0",
                status_code=400
            )
        if page_size > 30:
            return func.HttpResponse(
                "pageSize must be less than or equal to 30",
                status_code=400
            )
        
        
        continuation_token = req.params.get('continuationToken', None)
        if continuation_token is not None:
            try:
                base64.b64decode(continuation_token)
            except binascii.Error:
                return func.HttpResponse(
                    "Invalid continuationToken parameter",
                    status_code=400
                )
            
        
        min_square_meters = req.params.get('minSize', None)
        max_square_meters = req.params.get('maxSize', None)
        min_rent_price = req.params.get('minRentPrice', None)
        max_rent_price = req.params.get('maxRentPrice', None)
        listing_type = req.params.get('listingType', None)

        if listing_type is not None and listing_type not in ["rent", "sale"]:
            return func.HttpResponse(
                "Invalid listingType parameter",
                status_code=400
            )

        logging.info(f"Getting listings for city {partition_key} with page size {page_size}")

        query = "SELECT * FROM c WHERE c.partitionKey = @partitionKey"
        parameters = [
            { "name": "@partitionKey", "value": partition_key }
        ]

        if min_square_meters is not None:
            try:
                min_square_meters = int(min_square_meters)
            except ValueError:
                return func.HttpResponse(
                    "Invalid minSize parameter",
                    status_code=400
                )
            query += " AND c.squareMeters >= @minSize"
            parameters.append({ "name": "@minSize", "value": min_square_meters })

        if max_square_meters is not None:
            try:
                max_square_meters = int(max_square_meters)
            except ValueError:
                return func.HttpResponse(
                    "Invalid maxSize parameter",
                    status_code=400
                )
            query += " AND c.squareMeters <= @maxSize"
            parameters.append({ "name": "@maxSize", "value": max_square_meters })

        if min_rent_price is not None:
            try:
                min_rent_price = int(min_rent_price)
            except ValueError:
                return func.HttpResponse(
                    "Invalid minRentPrice parameter",
                    status_code=400
                )
            query += " AND c.rentPricePerMonth >= @minRentPrice"
            parameters.append({ "name": "@minRentPrice", "value": min_rent_price })

        if max_rent_price is not None:
            try:
                max_rent_price = int(max_rent_price)
            except ValueError:
                return func.HttpResponse(
                    "Invalid maxRentPrice parameter",
                    status_code=400
                )
            query += " AND c.rentPricePerMonth <= @maxRentPrice"
            parameters.append({ "name": "@maxRentPrice", "value": max_rent_price })

        if listing_type is not None:
            query += " AND (c.listingType = @listingType OR c.listingType = 'both')"
            parameters.append({ "name": "@listingType", "value": listing_type })

        container = client.get_database_client(database_name).get_container_client(container_name)

        query_iterable = container.query_items(
            query=query,
            parameters=parameters,
            partition_key=partition_key,
            max_item_count=page_size
        )

        try:
            pager = query_iterable.by_page(continuation_token)
            results = list(pager.next())
            continuation_token = pager.continuation_token
        except StopIteration:
            logging.warning("Failed to query listings: StopIteration")
            #return func.HttpResponse(
            #    "Failed to query listings: StopIteration",
            #    status_code=500
            #)
            results = []
            continuation_token = None
        except exceptions.CosmosHttpResponseError as e:
            # Probably because of malformed continuation token, but won't bother to check for now
            logging.error(f"Failed to query listings: {e}")
            return func.HttpResponse(
                "Failed to query listings",
                status_code=500
            )

        for result in results:
            for key in ["_rid", "_self", "_etag", "_attachments", "_ts"]:
                if key in result:
                    del result[key]

        response_body = {
            "results": results,
            "continuationToken": continuation_token
        }

        return func.HttpResponse(
            json.dumps(response_body),
            status_code=200,
            mimetype="application/json"
        )

    except Exception as e:
        logging.error(f"An error occurred: {e}")
        return func.HttpResponse(
            f"An error occurred",
            status_code=500
        )
