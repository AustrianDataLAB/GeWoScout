import azure.functions as func
import logging
import json
import os


bp = func.Blueprint() 


'''
NOTE Laurenz: 
- create a new module for each website inside scrapers/impl
- function name has to be unique: <genossenschaft_name>_scraper
- for debugging you can trigger the function via HTTP. Refer to https://learn.microsoft.com/en-us/azure/azure-functions/functions-manually-run-non-http?tabs=azure-portal
'''
@bp.timer_trigger(schedule="0 */5 * * * *", arg_name="timerObj", run_on_startup=False) 
@bp.queue_output(arg_name="q", queue_name=os.getenv('QUEUE_NAME'), connection="AzureWebJobsStorage")
def demo_scraper(timerObj: func.TimerRequest, q: func.Out[str]) -> None:
    logging.info('Scraper Demo A triggered.')
    
    payload = {
        "scraperId": "demoScraper",
        "timestamp": "2024-04-06T15:30:00Z",
        "listings": [
            {
                "title": "Modern 3-Bedroom Apartment in Central Vienna",
                "housingCooperative": "FutureLiving Genossenschaft",
                "projectId": "FLG2024",
                "listingId": "12345ABC",
                "country": "Austria",
                "city": "Vienna",
                "postalCode": "1010",
                "address": "Beispielgasse 42",
                "roomCount": 3,
                "squareMeters": 95,
                "availabilityDate": "2024-09-01",
                "yearBuilt": 2019,
                "hwgEnergyClass": "A",
                "fgeeEnergyClass": "A+",
                "listingType": "both",
                "rentPricePerMonth": 1200,
                "cooperativeShare": 5000,
                "salePrice": 350000,
                "additionalFees": 6500,
                "detailsUrl": "https://www.futurelivinggenossenschaft.at/listings/12345ABC",
                "previewImageUrl": "https://www.futurelivinggenossenschaft.at/listings/12345ABC/preview.jpg"
            }
        ]
    }

    q.set(json.dumps(payload).encode(encoding='UTF-8'))
