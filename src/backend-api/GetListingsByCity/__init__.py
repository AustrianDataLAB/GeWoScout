import json
import logging

import azure.functions as func


def main(req: func.HttpRequest, existingDoc: func.DocumentList) -> func.HttpResponse:
    results = []

    for doc in existingDoc:
        data = doc.to_dict()
        for key in ["_rid", "_self", "_ts", "_etag"]:
            data.pop(key, None)
        results.append(data)
    
    return func.HttpResponse(json.dumps(results), status_code=200, mimetype="application/json")
