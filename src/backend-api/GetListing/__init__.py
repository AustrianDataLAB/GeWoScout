import json
import logging

import azure.functions as func


def main(req: func.HttpRequest, existingDoc: func.DocumentList) -> func.HttpResponse:
    pk = req.route_params.get('city')
    id = req.route_params.get('id')
    logging.info('Retrieving listing %s : %s', pk, id)

    if not existingDoc or existingDoc.data == None or len(existingDoc.data) == 0:
        return func.HttpResponse(status_code=404)
    
    if len(existingDoc.data) > 1:
        logging.warning("Found more than one listing with the same id %s", id)

    data = existingDoc.data[0].to_dict()

    # Remove CosmosDB metadata
    for key in ["_rid", "_self", "_ts", "_etag"]:
        data.pop(key, None)

    return func.HttpResponse(json.dumps(data), status_code=200, mimetype="application/json")
