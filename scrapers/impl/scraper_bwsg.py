import azure.functions as func
from bs4 import BeautifulSoup
from datetime import datetime, timezone
import requests
import logging
import re
import json
import os

URL = "https://immobilien.bwsg.at"
PARAMS = {
    # 'f[all][marketing_type]': 'rent', # Miete
	'f[all][realty_type][0]': '2',    # Wohnung
    'f[all][realty_type][1]': '3',    # Haus
    'f[all][city]'          : 'Wien',
	'from'                  : '1117350'
}
QUEUE_BATCH_SIZE = 20

bp = func.Blueprint() 

@bp.timer_trigger(schedule="0 */5 * * * *", arg_name="timerObj", run_on_startup=False) 
@bp.queue_output(arg_name="q", queue_name=os.getenv('QUEUE_NAME'), connection="QUEUE_CONNECTION_STRING")
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

    payload = dict()
    payload["scraperId"] = 'bwsg_scraper_rent'

    
    payload["timestamp"] = datetime.now(timezone.utc).isoformat()
    Wohnungen = list()


    for i, link in enumerate(links):
        req = requests.request(method = 'GET', url = URL + link)
        soup = BeautifulSoup(req.text, 'html.parser')
        info = soup.find(class_ = 'container-wrapper')
        
        wohnung = dict()        
        headline = soup.find(class_ = 'realty-detail-headline')
        wohnung["title"] = headline.find('h1').get_text().strip()
        address = headline.find(class_ = 'address').get_text().strip()
        postal_code = address.split('\n')[0]
        postal_code = re.findall(r'\d+', postal_code)
        address = ''.join(address.split('\n')[1:])
        address = ' '.join(address.split())
        wohnung["address"] = address.strip()
        wohnung["postalCode"] = postal_code[0]
        wohnung["housingCooperative"] = "BWS-Gruppe"
        wohnung["country"] = "Austria"
        wohnung["city"] = "Vienna"
        
        detail_infos  = info.find(class_ = 'realty-detail-info').find_all('li')
        
        items = dict()
        for detail in detail_infos:
            desc = detail.find(class_ = 'list-item-desc').get_text().strip()
            value = detail.find(class_ = 'list-item-value').get_text().strip()
            items[desc] = value
        
        Objektnr = items.get("Objektnr.", "").split("/")
        
        if len(Objektnr) > 1:
            wohnung["projectId"] = Objektnr[0]
            wohnung["listingId"] = Objektnr[1]
        else:
            wohnung["projectId"] = Objektnr[0]
            wohnung["listingId"] = ""
            
        roomCount = items.get("Zimmer", None)
        wohnung["roomCount"] = int(roomCount) if roomCount is not None else None
        square_meters = re.findall(r'(\d+[\.\,\d]*)', items["Wohnfläche"])
        wohnung["squareMeters"] = float(square_meters[0].replace(',', '.'))
        wohnung["availabilityDate"] = items.get("Beziehbar", "")
        wohnung["yearBuilt"] = items.get("Baujahr", "")
        
        hwb = re.findall(r'(\d+[\.\,\d]*)', items.get("HWB", ""))
        # Pay attention here, the HWB uses a dot as a decimal separator
        hwb = float(hwb[0].replace(',', '.')) if len(hwb) > 0 else None
        
        hwb_class = re.findall(r'\b[A-G]{1}\b\+*', items.get("HWB", ""))
        hwb_class = hwb_class[0] if len(hwb_class) > 0 else None
        
        wohnung["hwgEnergyClass"] = hwb_class
        wohnung["hwgEnergy"] = hwb
        
        fgee = re.findall(r'(\d+[\.\,\d]*)', items.get("fGEE", ""))
        fgee = float(fgee[0].replace(',', '.')) if len(fgee) > 0 else None
        
        fgee_class = re.findall(r'\b[A-G]{1}\b\+*', items.get("fGEE", ""))
        fgee_class = fgee_class[0] if len(fgee_class) > 0 else None
        
        wohnung["fgeeEnergyClass"] = fgee_class
        wohnung["fgeeEnergy"] = fgee
        
        detail_preis = info.find(class_ = 'realty-detail-prices')

        rent_preis  = detail_preis.find(class_ = 'rent-price-table').find_all('tr')
        
        rent = dict()
        for row in rent_preis:
            cols = row.find_all('td')
            desc = cols[0].get_text().strip()
            value = cols[1].get_text().strip()
            rent[desc] = value
            
        miete = rent.get("Miete:", "")
        miete = re.findall(r'(\d+[\.\,\d]*)', miete)
        miete = float(miete[0].replace('.', '')\
            .replace(',', '.')) if len(miete) > 0 else None
        
        wohnung["rentPricePerMonth"] = miete
        
        kaufpreis = rent.get("Kaufpreis:", "")
        kaufpreis = re.findall(r'(\d+[\.\,\d]*)', kaufpreis)
        kaufpreis = float(kaufpreis[0].replace('.', '')\
            .replace(',', '.')) if len(kaufpreis) > 0 else None
        
        wohnung["salePrice"] = kaufpreis
        
        add_rent = rent.get("monatliche Gesamtbelastung:", "")
        add_rent = re.findall(r'(\d+[\.\,\d]*)', add_rent)
        add_rent = float(add_rent[0].replace('.', '')\
            .replace(',', '.')) if len(add_rent) > 0 else None
        
        if miete is not None:
            wohnung["additionalFees"] = add_rent - miete
        else:
            wohnung["additionalFees"] = add_rent
            
        extras = detail_preis.find(class_ = 'list-unstyled').find_all('li')
        
        extras_cost = dict()
        for extra in extras:
            desc = extra.find(class_ = 'list-item-desc').get_text().strip()
            value = extra.find(class_ = 'list-item-value').get_text().strip()
            extras_cost[desc] = value
        
        deposit = extras_cost.get("Kaution:", "")
        deposit = re.findall(r'(\d+[\.\,\d]*)', deposit)
        deposit = float(deposit[0].replace('.', '')\
            .replace(',', '.')) if len(deposit) > 0 else None
        
        wohnung["deposit"] = deposit
        
        fin_contr = extras_cost.get("Finanzierungsbeitrag:", "")
        fin_contr = re.findall(r'\d+\,*\.*\d*\,*\d*', fin_contr)
        fin_contr = float(fin_contr[0].replace('.', '')\
            .replace(',', '.')) if len(fin_contr) > 0 else None
        
        wohnung["financialContribution"] = fin_contr
        
        coop_share = extras_cost.get("Provision:", "")
        coop_share = re.findall(r'\d+\,*\.*\d*\,*\d*', coop_share)
        coop_share = float(coop_share[0].replace('.', '')\
            .replace(',', '.')) if len(coop_share) > 0 else None
        
        wohnung["cooperativeShare"] = coop_share
        
        if miete is not None and kaufpreis is not None:
            wohnung["listingType"] = "Both"
        elif kaufpreis is not None:
            wohnung["listingType"] = "Sale"
        else:
            wohnung["listingType"] = "Rent"
            
        wohnung["detailsUrl"] = URL + link
        
        images = soup.find_all(class_ = 'carousel-inner')
        
        
        # Returning the first image
        image = images[0].find_all('img')[0].get('src')
        
        wohnung["previewImageUrl"] = image
        
        
        explanation = info.find(class_ = 'costs-explanation')
        
        if explanation is not None:
            explanation = explanation.get_text().strip()
        
        
        wohnung['cost-explanation'] = explanation
        
        Wohnungen.append(wohnung)

        if i % QUEUE_BATCH_SIZE == 0:
            payload["listings"] = Wohnungen
            q.set(json.dumps(payload).encode(encoding='UTF-8'))
            del payload["listings"]
            Wohnungen = list()
    
    payload["listings"] = Wohnungen
    
    q.set(json.dumps(payload).encode(encoding='UTF-8'))
