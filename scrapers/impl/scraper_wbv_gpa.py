import json
import logging
import os
import re
from datetime import timezone, datetime
from typing import Optional

import azure.functions as func
import requests
from bs4 import BeautifulSoup

URL = "https://www.wbv-gpa.at/wohnungen/"

bp = func.Blueprint()


def extract_postal_code_and_street(address: str) -> (str, str):
    """
    Extracts postal code and street address from a string.

    Args:
        address (str): Address string.

    Returns:
        tuple: Postal code and street address.
    """
    postal_code = re.findall(r'\d{4}', address)[0]
    street_address = re.sub(r'\d{4},?\s*Wien,?\s*', '', address).strip()
    return postal_code, street_address


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
    vienna_listings_items = soup.find_all(class_="wien")
    for item in vienna_listings_items:
        try:
            link = item.find('a').get('href')
            postal_code, street_address = extract_postal_code_and_street(item.attrs['data-location'])

            listing = {
                "yearBuilt": int(item.attrs['data-year']),
                "squareMeters": float(item.attrs['data-space'].replace(',', '.')),
                "detailsUrl": link,
                "city": "Vienna",
                "country": "Austria",
                "housingCooperative": "WBV-GPA",
                "projectId": "",
                "listingId": "",
                "postalCode": postal_code,
                "address": street_address,
                "availabilityDate": "sofort"
            }

            if "data-rent" in item.attrs:
                listing["listingType"] = "rent"
                listing["rentPricePerMonth"] = float(item.attrs['data-rent'].replace(',', '.'))
                listing["cooperativeShare"] = float(item.attrs['data-financing'].replace(".", "").replace(',', '.'))

            listings.append(listing)
        except Exception as e:
            logging.error("Error while extracting listing info from url %s", e)

    return listings


def extract_room_number(value: str) -> int:
    """
    This function is used to extract the room number from a given string.

    Parameters:
        value (str): The string from which to extract the room number.

    Returns:
        int: The room number extracted from the string.
    """
    return int(re.findall(r'\d+', value)[0])


def extract_energy_class(value: str) -> (str, str):
    """
    This function is used to extract energy class information from a given string.

    Parameters:
        value (str): The string from which to extract the energy class information.

    Returns:
        tuple: A tuple containing two elements. The first element is the energy class type.
        The second element is the energy class value (A-G). Returns (None, None) if no energy class is found.
    """

    # Find all occurrences of energy classes (A-G) in the input string
    hwb_class = re.findall(r'\b[A-G]\b\+*', value)

    # If no energy class is found, return None, None
    if len(hwb_class) == 0:
        return None, None

    # If "hwb" is found in the input string, return "hwb" and the first found energy class
    if "hwb" in value.lower():
        return "hwb", hwb_class[0]

    # If "fgee" is found in the input string, return "fgee" and the first found energy class
    if "fgee" in value.lower():
        return "fgee", hwb_class[0]

    # If neither "hwb" nor "fgee" is found, return None, None
    return None, None


def add_details_to_listing(listing: dict, listing_soup: BeautifulSoup) -> Optional[dict]:
    """
    Adds additional details from the listing page to the listing dictionary.

    Args:
        listing (dict): The listing dictionary to which details should be added.
        listing_soup (BeautifulSoup): BeautifulSoup object containing the HTML of the listing page.

    Returns:
        dict: The updated listing dictionary with additional details.
        If no title is found for the listing, the function returns None and logs an info message.
    """

    head_line = listing_soup.find(class_=["hero__content__title", "h1"])
    if head_line is None:
        logging.info("No title found for listing %s, skipping", listing['detailsUrl'])
        return None

    listing["title"] = head_line.get_text().strip()

    for detail_soup in listing_soup.find_all(class_="projectsingle__details__item__title"):
        key = detail_soup.get_text().strip()
        value = detail_soup.find_next_sibling().get_text().strip()

        if key == "Einheiten":
            listing['roomCount'] = extract_room_number(value)
        elif key == "Energie":
            energy_class, energy_value = extract_energy_class(value)
            if energy_class == "hwb":
                listing['hwgEnergyClass'] = energy_value
            elif energy_class == "fgee":
                listing['fgeeEnergyClass'] = energy_value

    for image_soup in listing_soup.find_all(class_="slideshow__inner__container__image"):
        img_url = image_soup.find("img").get("src")
        listing["previewImageUrl"] = img_url
        break

    return listing


@bp.timer_trigger(schedule="0 */5 * * * *", arg_name="timerObj", run_on_startup=False)
@bp.queue_output(arg_name="q", queue_name=os.getenv('QUEUE_NAME'), connection="AzureWebJobsStorage")
def wbv_gpa_scraper(timerObj: func.TimerRequest, q: func.Out[str]) -> None:
    logging.info('WBV-GPA scraper triggered.')

    req = requests.request(method='GET', url=URL)
    soup = BeautifulSoup(req.text, 'html.parser')

    try:
        vienna_listing_links = extract_listing_info(soup)

        output = {
            "scraperId": "wbv_gpa_scraper",
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "listings": []
        }
        for listing in vienna_listing_links:
            req = requests.request(method='GET', url=listing['detailsUrl'])
            listing_soup = BeautifulSoup(req.text, 'html.parser')

            if (updated_listing := add_details_to_listing(listing, listing_soup)) is None:
                continue

            output['listings'].append(updated_listing)
            if len(output['listings']) >= 20:
                q.set(json.dumps(output).encode(encoding='UTF-8'))
                del output['listings']
                output['listings'] = []

        if len(output['listings']) > 0:
            q.set(json.dumps(output).encode(encoding='UTF-8'))
        logging.info('WBV-GPA scraper finished.')
    except Exception as e:
        logging.error("Error while scraping WBV-GPA: %s", e)

