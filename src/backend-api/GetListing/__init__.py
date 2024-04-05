import json
import logging

import azure.functions as func


def main(req: func.HttpRequest, existingDoc: func.DocumentList) -> func.HttpResponse:
    id = req.route_params.get('id')
    logging.info('Retrieving listing %s', id)

    if not existingDoc or existingDoc.data == None or len(existingDoc.data) == 0:
        return func.HttpResponse("Listing not found", status_code=404)
    
    if len(existingDoc.data) > 1:
        logging.warn("Found more than one listing with the same id %s", id)

    data = existingDoc.data[0].to_dict()

    # Remove CosmosDB metadata
    if "_rid" in data:
        del data["_rid"]
    if "_self" in data:
        del data["_self"]
    if "_ts":
        del data["_ts"]
    if "_etag":
        del data["_etag"]

    return func.HttpResponse(json.dumps(data), status_code=200, mimetype="application/json")
