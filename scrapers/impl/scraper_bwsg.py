import json
import logging
import os
import re
from datetime import datetime, timezone
from typing import Optional

import azure.functions as func
import requests
from bs4 import BeautifulSoup

URL = "https://immobilien.bwsg.at"
PARAMS = {
    # 'f[all][marketing_type]': 'rent', # Miete
    'f[all][realty_type][0]': '2',  # Wohnung
    'f[all][realty_type][1]': '3',  # Haus
    'f[all][city]': 'Wien',
    'from': '1117350'
}
QUEUE_BATCH_SIZE = 20

bp = func.Blueprint()

def format_date(date: str):
    # TODO implement other common BWSG patterns
    if re.match(r"\d{1,2}\.\d{1,2}\.\d{4}", date):
        date_obj = datetime.strptime(date, "%d.%m.%Y")
        return date_obj.strftime("%Y-%m-%d")
    else:
        # if date could not be parsed lets assume it's available from today
        return datetime.now().strftime("%Y-%m-%d")


def extract_energy_class(raw_data: dict, field_name: str) -> Optional[str]:
    class_extraction = re.findall(r'\b[A-G]\b\+*', raw_data.get(field_name, ""))
    class_value = class_extraction[0] if len(class_extraction) > 0 else None
    return class_value


@bp.timer_trigger(schedule="0 */5 * * * *", arg_name="timerObj", run_on_startup=False)
@bp.queue_output(arg_name="q", queue_name=os.getenv('QUEUE_NAME'), connection="QUEUE_CONNECTION_STRING")
def bwsg_scraper(timerObj: func.TimerRequest, q: func.Out[str]) -> None:
    logging.info('BWSG scraper triggered.')

    req = requests.request(method='GET', url=URL, params=PARAMS)
    soup = BeautifulSoup(req.text, 'html.parser')
    pages = soup.find(class_='pagination')
    cur = pages.find(class_='active')
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
        pages = soup.find(class_='pagination')
        cur = pages.find(class_='active')
        pages = pages.find_all('li')

    payload = dict()
    payload["scraperId"] = 'bwsg_scraper_rent'

    payload["timestamp"] = datetime.now(timezone.utc).isoformat()
    scraped_listings = list()

    for i, link in enumerate(links):
        req = requests.request(method='GET', url=URL + link)
        soup = BeautifulSoup(req.text, 'html.parser')
        info = soup.find(class_='container-wrapper')

        listing = dict()
        headline = soup.find(class_='realty-detail-headline')
        listing["title"] = headline.find('h1').get_text().strip()
        address = headline.find(class_='address').get_text().strip()
        postal_code = address.split('\n')[0]
        postal_code = re.findall(r'\d+', postal_code)
        address = ''.join(address.split('\n')[1:])
        address = ' '.join(address.split())
        listing["address"] = address.strip()
        listing["postalCode"] = postal_code[0]
        listing["housingCooperative"] = "BWS-Gruppe"
        listing["country"] = "Austria"
        listing["city"] = "Vienna"

        detail_infos = info.find(class_='realty-detail-info').find_all('li')

        bwsg_description = dict()
        for detail in detail_infos:
            desc = detail.find(class_='list-item-desc').get_text().strip()
            value = detail.find(class_='list-item-value').get_text().strip()
            bwsg_description[desc] = value

        object_nr = bwsg_description.get("Objektnr.", "").split("/")

        if len(object_nr) > 1:
            listing["projectId"] = object_nr[0]
            listing["listingId"] = object_nr[1]
        else:
            listing["projectId"] = object_nr[0]
            listing["listingId"] = ""

        room_count = bwsg_description.get("Zimmer", None)
        listing["roomCount"] = int(room_count) if room_count is not None else None
        if "Wohnfläche" in bwsg_description:
            square_meters = re.findall(r'(\d+[.,\d]*)', bwsg_description["Wohnfläche"])
        listing["squareMeters"] = float(square_meters[0].replace(',', '.'))
        listing["availabilityDate"] = format_date(bwsg_description.get("Beziehbar", ""))
        listing["yearBuilt"] = int(bwsg_description.get("Baujahr", "-1"))

        if (hwb_energy_class := extract_energy_class(bwsg_description, "HWB")) is not None:
            listing["hwgEnergyClass"] = hwb_energy_class
        if (fgee_energy_class := extract_energy_class(bwsg_description, "fGEE")) is not None:
            listing["fgeeEnergyClass"] = fgee_energy_class

        detail_preis = info.find(class_='realty-detail-prices')
        rent_preis = detail_preis.find(class_='rent-price-table').find_all('tr')

        rent = dict()
        for row in rent_preis:
            cols = row.find_all('td')
            desc = cols[0].get_text().strip()
            value = cols[1].get_text().strip()
            rent[desc] = value

        miete = rent.get("Miete:", "")
        miete = re.findall(r'(\d+[\.\,\d]*)', miete)
        miete = float(miete[0].replace('.', '') \
                      .replace(',', '.')) if len(miete) > 0 else None

        listing["rentPricePerMonth"] = miete

        kaufpreis = rent.get("Kaufpreis:", "")
        kaufpreis = re.findall(r'(\d+[\.\,\d]*)', kaufpreis)
        kaufpreis = float(kaufpreis[0].replace('.', '') \
                          .replace(',', '.')) if len(kaufpreis) > 0 else None

        listing["salePrice"] = kaufpreis

        add_rent = rent.get("monatliche Gesamtbelastung:", "")
        add_rent = re.findall(r'(\d+[\.\,\d]*)', add_rent)
        add_rent = float(add_rent[0].replace('.', '') \
                         .replace(',', '.')) if len(add_rent) > 0 else None

        if miete is not None:
            listing["additionalFees"] = add_rent - miete
        else:
            listing["additionalFees"] = add_rent

        extras = detail_preis.find(class_='list-unstyled').find_all('li')

        extras_cost = dict()
        for extra in extras:
            desc = extra.find(class_='list-item-desc').get_text().strip()
            value = extra.find(class_='list-item-value').get_text().strip()
            extras_cost[desc] = value

        deposit = extras_cost.get("Kaution:", "")
        deposit = re.findall(r'(\d+[.,\d]*)', deposit)
        deposit = float(deposit[0].replace('.', '') \
                        .replace(',', '.')) if len(deposit) > 0 else None

        listing["deposit"] = deposit

        fin_contr = extras_cost.get("Finanzierungsbeitrag:", "")
        fin_contr = re.findall(r'\d+,*\.*\d*,*\d*', fin_contr)
        fin_contr = float(fin_contr[0].replace('.', '') \
                          .replace(',', '.')) if len(fin_contr) > 0 else None

        listing["financialContribution"] = fin_contr

        coop_share = extras_cost.get("Provision:", "")
        coop_share = re.findall(r'\d+,*\.*\d*,*\d*', coop_share)
        coop_share = float(coop_share[0].replace('.', '') \
                           .replace(',', '.')) if len(coop_share) > 0 else None

        listing["cooperativeShare"] = coop_share

        if miete is not None and kaufpreis is not None:
            listing["listingType"] = "both"
        elif kaufpreis is not None:
            listing["listingType"] = "sale"
        else:
            listing["listingType"] = "rent"

        listing["detailsUrl"] = URL + link

        images = soup.find_all(class_='carousel-inner')

        # Returning the first image
        image = images[0].find_all('img')[0].get('src')

        listing["previewImageUrl"] = image

        explanation = info.find(class_='costs-explanation')

        if explanation is not None:
            explanation = explanation.get_text().strip()

        listing['cost-explanation'] = explanation

        scraped_listings.append(listing)

        if i % QUEUE_BATCH_SIZE == 0:
            payload["listings"] = scraped_listings
            q.set(json.dumps(payload).encode(encoding='UTF-8'))
            del payload["listings"]
            scraped_listings = list()

    payload["listings"] = scraped_listings

    q.set(json.dumps(payload).encode(encoding='UTF-8'))
