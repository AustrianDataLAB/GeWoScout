import unittest
from unittest.mock import patch

from bs4 import BeautifulSoup

from impl.scraper_wbv_gpa import extract_energy_class, extract_room_number, extract_postal_code_and_street, \
    extract_listing_info, add_details_to_listing


class TestExtractPostalCodeAndStreet(unittest.TestCase):

    def test_extract_postal_code_and_street(self):
        address = "1234, Wien, Some Street"
        postal_code, street_address = extract_postal_code_and_street(address)
        self.assertEqual(postal_code, "1234")
        self.assertEqual(street_address, "Some Street")


class TestExtractListingInfo(unittest.TestCase):

    def setUp(self):
        self.html_doc = """
        <div class="wien" data-location="1234, Wien, Some Street" data-year="2000" data-space="100,5" data-rent="500,00" data-financing="1000,00">
            <a href="http://example.com">Link</a>
        </div>
        """
        self.soup = BeautifulSoup(self.html_doc, 'html.parser')

    def test_extract_listing_info(self):
        listings = extract_listing_info(self.soup)
        self.assertEqual(len(listings), 1)
        self.assertEqual(listings[0]['yearBuilt'], 2000)
        self.assertEqual(listings[0]['squareMeters'], 100.5)
        self.assertEqual(listings[0]['detailsUrl'], 'http://example.com')
        self.assertEqual(listings[0]['postalCode'], '1234')
        self.assertEqual(listings[0]['address'], 'Some Street')
        self.assertEqual(listings[0]['listingType'], 'rent')
        self.assertEqual(listings[0]['rentPricePerMonth'], 500.00)
        self.assertEqual(listings[0]['cooperativeShare'], 1000.00)

    @patch('impl.scraper_wbv_gpa.logging.error')
    def test_extract_listing_info_error(self, mock_logging_error):
        html_doc = """
        <div class="wien" data-location="wrong address format">
            <a href="http://example.com">Link</a>
        </div>
        """
        soup = BeautifulSoup(html_doc, 'html.parser')
        listings = extract_listing_info(soup)
        self.assertEqual(len(listings), 0)
        mock_logging_error.assert_called_once()


class TestExtractRoomNumber(unittest.TestCase):

    def test_extract_room_number(self):
        value = "3 rooms"
        self.assertEqual(extract_room_number(value), 3)

    def test_extract_room_number_german(self):
        value = "3 Zimmer"
        self.assertEqual(extract_room_number(value), 3)


class TestExtractEnergyClass(unittest.TestCase):
    def test_extract_energy_class_hwb(self):
        value = "hwb 30 A+"
        energy_class, energy_value = extract_energy_class(value)
        self.assertEqual(energy_class, "hwb")
        self.assertEqual(energy_value, "A+")

    def test_extract_energy_class_fgee(self):
        value = "fgee 30 B"
        energy_class, energy_value = extract_energy_class(value)
        self.assertEqual(energy_class, "fgee")
        self.assertEqual(energy_value, "B")

    def test_extract_energy_class_none(self):
        value = "no energy info"
        energy_class, energy_value = extract_energy_class(value)
        self.assertIsNone(energy_class)
        self.assertIsNone(energy_value)


class TestListingDetailsScraping(unittest.TestCase):

    @patch('logging.info')
    def test_add_details_to_listing_no_title(self, mock_logging):
        listing = {"detailsUrl": "test_url"}
        listing_soup = BeautifulSoup('<html></html>', 'html.parser')
        result = add_details_to_listing(listing, listing_soup)
        self.assertIsNone(result)
        mock_logging.assert_called_once_with("No title found for listing %s, skipping", "test_url")

    @patch('logging.info')
    def test_add_details_to_listing_with_title(self, mock_logging):
        listing = {"detailsUrl": "test_url"}
        html = """
        <html>
            <div class="hero__content__title h1">Test Title</div>
        </html>
        """
        listing_soup = BeautifulSoup(html, 'html.parser')
        result = add_details_to_listing(listing, listing_soup)
        self.assertIsNotNone(result)
        self.assertEqual(result["title"], "Test Title")

    def test_add_details_to_listing_with_room_count(self):
        listing = {"detailsUrl": "test_url"}
        html = """
        <html>
            <div class="hero__content__title h1">Test Title</div>
            <div class="projectsingle__details__item__title">Einheiten</div>
            <div>3 rooms</div>
        </html>
        """
        listing_soup = BeautifulSoup(html, 'html.parser')
        result = add_details_to_listing(listing, listing_soup)
        self.assertEqual(result["roomCount"], 3)

    def test_add_details_to_listing_with_hwb_energy_class(self):
        listing = {"detailsUrl": "test_url"}
        html = """
        <html>
            <div class="hero__content__title h1">Test Title</div>
            <div class="projectsingle__details__item__title">Energie</div>
            <div>hwb B</div>
        </html>
        """
        listing_soup = BeautifulSoup(html, 'html.parser')
        result = add_details_to_listing(listing, listing_soup)
        self.assertEqual(result["hwgEnergyClass"], "B")

    def test_add_details_to_listing_with_fgee_energy_class(self):
        listing = {"detailsUrl": "test_url"}
        html = """
        <html>
            <div class="hero__content__title h1">Test Title</div>
            <div class="projectsingle__details__item__title">Energie</div>
            <div>fgEe B</div>
        </html>
        """
        listing_soup = BeautifulSoup(html, 'html.parser')
        result = add_details_to_listing(listing, listing_soup)
        self.assertEqual(result["fgeeEnergyClass"], "B")

    def test_add_details_to_listing_with_image_url(self):
        listing = {"detailsUrl": "test_url"}
        html = """
        <html>
            <div class="hero__content__title h1">Test Title</div>
            <div class="slideshow__inner__container__image">
                <img src="test_image_url">
            </div>
        </html>
        """
        listing_soup = BeautifulSoup(html, 'html.parser')
        result = add_details_to_listing(listing, listing_soup)
        self.assertEqual(result["previewImageUrl"], "test_image_url")


if __name__ == '__main__':
    unittest.main()
