import json
import logging
import os
import re
from datetime import timezone, datetime
from typing import Optional

import azure.functions as func
import requests
from bs4 import BeautifulSoup

URL = "https://www.oevw.at"

bp = func.Blueprint()


def extract_listing_info(soup: BeautifulSoup):
    """
    Extracts listing information from a BeautifulSoup object.

    Args:
        soup (BeautifulSoup): BeautifulSoup object containing the HTML of a webpage.

    Returns:
        list: A list of dictionaries, each representing a listing. If an error occurs while extracting
        information for a listing, an error message is logged and the listing is skipped.
    """
    listings = []
    container = soup.find_all(id="searchresults").pop()
    #vienna_listings_items = container.find_all_next(class_="thumb__content")
    vienna_listings_items = container.find_all_next(class_="thumb thumb--unit imagezoom-delegate-hover")
    for item in vienna_listings_items:
        try:
            type = item.find(class_="thumb__info small").get_text()
            # discard Geschäftslokale usw.
            if not "Wohnung" in type:
                continue

            postal_code = type.split("–")[-1].replace("Wien", "").strip()
            link = f"{URL}{item.find('a').get('href')}"
            title =item.find(class_="thumb__heading big").get_text().strip()
            square_meters = 0
            rent = 0
            room_count = 0
            fee = 0
            preview_url = f"{URL}{item.find('img').get('src')}"
            project_id = item.find(class_="thumb__project small").get_text()

            if "Miete mit Kaufoption" in type:
                type = "both"
            else:
                type = "rent"

            for detail in item.find(class_="thumb__subheading").find("ul").find_all("li"):
                text = detail.get_text()
                if "m²" in text:
                    square_meters = float(text.replace("m²", "").replace(",", ".").strip())

                if "€" in text:
                    rent = float(text.replace("€", "").replace(",", ".").strip())

            for detail in item.find(class_="thumb__text").find("ul").find_all("li"):
                text = detail.get_text()
                if "Zimmer" in text:
                    room_count = int(text.replace("Zimmer", "").strip())

                if "Eigenmittel" in text:
                    text = text.replace("Eigenmittel: €", "").replace(",", ".").strip()
                    if text.count('.') > 1:
                        first_period_index = text.index('.')
                        # Keep the part of the string up to and including the first period, then remove all other periods
                        text = text[:first_period_index + 1] + text[first_period_index + 1:].replace('.', '')

                    fee = float(text)

            listing = {
                "title": title,
                "housingCooperative": "OEVW",
                "projectId": project_id,
                "listingId": "",
                "country": "Austria",
                "city": "Vienna",
                "postalCode": postal_code,
                "address": "",
                "roomCount": room_count,
                "squareMeters": square_meters,
                "availabilityDate": datetime.now().strftime("%Y-%m-%d"),
                "yearBuilt": 0, # TODO: remove
                "listingType": type,
                "detailsUrl": link,
                "previewImageUrl": preview_url,
                "rentPricePerMonth": rent,
                "cooperativeShare": fee,
            }

            listings.append(listing)
        except Exception as e:
            logging.error("Error while extracting listing info from url %s", e)

    return listings


def add_details_to_listing(listing: dict, listing_soup: BeautifulSoup) -> Optional[dict]:
    """
    Adds additional details (address) from the listing page to the listing dictionary.

    Args:
        listing (dict): The listing dictionary to which details should be added.
        listing_soup (BeautifulSoup): BeautifulSoup object containing the HTML of the listing page.

    Returns:
        dict: The updated listing dictionary with additional details.
        If no address is found for the listing, the function returns the listing without the address.
    """

    try:
        address = listing_soup.find(class_="description__address big").get_text().split(",")[0]
        listing["address"] = address
    except Exception as e:
        logging.error("Error while extracting listing info from url %s", e)
        return listing

    return listing


@bp.timer_trigger(schedule="0 */5 * * * *", arg_name="timerObj", run_on_startup=False)
@bp.queue_output(arg_name="q", queue_name=os.getenv('QUEUE_NAME'), connection="AzureWebJobsStorage")
def oevw_scraper(timerObj: func.TimerRequest, q: func.Out[str]) -> None:
    logging.info('OEVW scraper triggered.')

    # each page contains 12 listings
    page = 1
    while True:
        req_url = f"{URL}/suche?page={page}"
        req = requests.request(method='GET', url=req_url)

        # break if automatically redirected to page 1
        if req.url != req_url:
            break

        soup = BeautifulSoup(req.text, 'html.parser')

        try:
            vienna_listing_links = extract_listing_info(soup)

            output = {
                "scraperId": "oevw_scraper",
                "timestamp": datetime.now(timezone.utc).isoformat(),
                "listings": []
            }
            for listing in vienna_listing_links:
                req = requests.request(method='GET', url=listing['detailsUrl'])
                listing_soup = BeautifulSoup(req.text, 'html.parser')

                if (updated_listing := add_details_to_listing(listing, listing_soup)) is None:
                    continue

                output['listings'].append(updated_listing)

            if len(output['listings']) > 0:
                q.set(json.dumps(output).encode(encoding='UTF-8'))
                del output['listings']
                output['listings'] = []

        except Exception as e:
            logging.error("Error while scraping OEVW: %s", e)

    logging.info('OEVW scraper finished.')
