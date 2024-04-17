import azure.functions as func
import datetime
import json
import logging

app = func.FunctionApp()

# NOTE (laurenz): I left this HTTP triggered function in for easier debugging - remove for prod
@app.route(route="http_trigger", auth_level=func.AuthLevel.ANONYMOUS)
def http_trigger(req: func.HttpRequest) -> func.HttpResponse:
    logging.info('Python HTTP trigger function processed a request.')

    name = req.params.get('name')

    if name:
        return func.HttpResponse(f"Hola {name}.")
    else:
        return func.HttpResponse("Grias di", status_code=200)



@app.timer_trigger(schedule="0 0 0 * * *", arg_name="myTimer", run_on_startup=False,
              use_monitor=False) 
def timer_trigger(myTimer: func.TimerRequest) -> None:
    
    if myTimer.past_due:
        logging.info('The timer is past due!')

    logging.info('Python timer trigger function executed.')