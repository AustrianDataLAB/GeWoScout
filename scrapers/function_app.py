import azure.functions as func

app = func.FunctionApp()

# NOTE (Laurenz): Register the created scrapers in the function app
from impl import scraper_bwsg
from impl import scraper_demo

app.register_functions(scraper_bwsg.bp)
app.register_functions(scraper_demo.bp)
