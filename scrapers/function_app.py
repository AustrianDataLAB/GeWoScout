import azure.functions as func

app = func.FunctionApp()

from impl import scraper_bwsg
from impl import scraper_demo

app.register_functions(scraper_bwsg.bp)
app.register_functions(scraper_demo.bp)
