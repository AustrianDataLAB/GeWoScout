import azure.functions as func

app = func.FunctionApp()

# NOTE (Laurenz): Register the created scrapers in the function app
from impl import scraper_bwsg
app.register_functions(scraper_bwsg.bp)
