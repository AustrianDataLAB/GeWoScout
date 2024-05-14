import azure.functions as func
from bs4 import BeautifulSoup
import requests
import logging
import json
import os

URL = "https://immobilien.bwsg.at"
PARAMS = {
    'f[all][marketing_type]': 'rent', # Miete
	'f[all][realty_type][0]': '2',    # Wohnung
    'f[all][realty_type][1]': '3',    # Haus
    'f[all][city]'          : 'Wien',
	'from'                  : '1117350'
}

bp = func.Blueprint() 

@bp.timer_trigger(schedule="0 */5 * * * *", arg_name="timerObj", run_on_startup=False) 
@bp.queue_output(arg_name="q", queue_name=os.getenv('QUEUE_NAME'), connection="AzureWebJobsStorage")
def bwsg_scraper(timerObj: func.TimerRequest, q: func.Out[str]) -> None:
    logging.info('BWSG scraper triggered.')

    req = requests.request(method='GET', url=URL, params=PARAMS)
    soup = BeautifulSoup(req.text, 'html.parser')
    pages = soup.find(class_ = 'pagination')
    cur = pages.find(class_ = 'active')
    pages = pages.find_all('li')

    links = []

    while True:
        panel_wrapper = soup.find_all(class_='panel-wrapper')

        for panel in panel_wrapper:
            panel_footer = panel.find(class_='panel-footer')
            links.append(panel.find('a').get('href'))
        
        if cur == pages[-1]:
            break
                
        cur_index = pages.index(cur)
        next_index = cur_index + 1
        next_link = pages[next_index].find('a').get('href')
        
        req = requests.request(method='GET', url=URL + next_link)
        soup = BeautifulSoup(req.text, 'html.parser')
        pages = soup.find(class_ = 'pagination')
        cur = pages.find(class_ = 'active')
        pages = pages.find_all('li')    

    Wohnungen = list()

    count = 0

    for link in links:        
        req = requests.request(method = 'GET', url = URL + links[0])
        soup = BeautifulSoup(req.text, 'html.parser')
        info = soup.find(class_ = 'container-wrapper')
        detail_infos  = info.find(class_ = 'realty-detail-info').find_all('li')
        
        wohnung = dict()
        for detail in detail_infos:
            desc = detail.find(class_ = 'list-item-desc').get_text().strip()
            value = detail.find(class_ = 'list-item-value').get_text().strip()
            wohnung[desc] = value
        
        detail_preis  = info.find(class_ = 'rent-price-table w-100').find_all('tr')
        
        for row in detail_preis:
            cols = row.find_all('td')
            desc = cols[0].get_text().strip()
            value = cols[1].get_text().strip()
            wohnung[desc] = value
            
        extra_info = info.find(class_ = 'list-unstyled').find_all('li')
        for detail in detail_infos:
            desc = detail.find(class_ = 'list-item-desc').get_text().strip()
            value = detail.find(class_ = 'list-item-value').get_text().strip()
            wohnung[desc] = value
        
        explanation = info.find(class_ = 'costs-explanation').get_text().strip()
        wohnung['cost-explanation'] = explanation
        
        wohnung['link'] = URL + links[0]
        wohnung['Unternehmen'] = 'BWSG'
        count += 1
        Wohnungen.append(wohnung)
    
    q.set(json.dumps(Wohnungen).encode(encoding='UTF-8'))
