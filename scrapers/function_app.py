import azure.functions as func

app = func.FunctionApp()

from impl import scraper_bwsg
from impl import scraper_wbv_gpa
from impl import scraper_oevw
# from impl import scraper_demo

app.register_functions(scraper_bwsg.bp)
app.register_functions(scraper_wbv_gpa.bp)
app.register_functions(scraper_oevw.bp)
# app.register_functions(scraper_demo.bp)
