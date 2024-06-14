package models

import "time"

type ScraperResultList struct {
	ScraperId string                 `json:"scraperId"`
	Timestamp time.Time              `json:"timestamp"`
	Listings  []ScraperResultListing `json:"listings"`
}

type ScraperResultListing struct {
	Title              string  `json:"title"`
	HousingCooperative string  `json:"housingCooperative"`
	ProjectId          string  `json:"projectId"`
	ListingId          string  `json:"listingId"`
	Country            string  `json:"country"`
	City               string  `json:"city"`
	PostalCode         string  `json:"postalCode"`
	Address            string  `json:"address"`
	RoomCount          int     `json:"roomCount"`
	SquareMeters       int     `json:"squareMeters"`
	AvailabilityDate   string  `json:"availabilityDate"`
	YearBuilt          int     `json:"yearBuilt"`
	HwgEnergyClass     *string `json:"hwgEnergyClass,omitempty"`
	FgeeEnergyClass    *string `json:"fgeeEnergyClass,omitempty"`
	ListingType        string  `json:"listingType"`
	RentPricePerMonth  *int    `json:"rentPricePerMonth,omitempty"`
	CooperativeShare   *int    `json:"cooperativeShare,omitempty"`
	SalePrice          *int    `json:"salePrice,omitempty"`
	AdditionalFees     *int    `json:"additionalFees,omitempty"`
	DetailsUrl         string  `json:"detailsUrl"`
	PreviewImageUrl    string  `json:"previewImageUrl"`
}

const ScraperResultListingSchema = `
	{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"title": "FlatListingMessage",
		"type": "object",
		"properties": {
			"scraperId": {
			"type": "string",
			"description": "A unique identifier for the scraper that fetched the data."
			},
			"timestamp": {
			"type": "string",
			"format": "date-time",
			"description": "The exact date and time when the scraping operation was performed."
			},
			"listings": {
			"type": "array",
			"items": {
				"type": "object",
				"properties": {
				"title": {
					"type": "string",
					"description": "The title or headline of the flat listing."
				},
				"housingCooperative": {
					"type": "string",
					"description": "The name of the housing cooperative or Genossenschaft."
				},
				"projectId": {
					"type": "string",
					"description": "The project identifier within the housing cooperative's portfolio."
				},
				"listingId": {
					"type": "string",
					"description": "A unique identifier for the individual flat listing."
				},
				"country": {
					"type": "string",
					"description": "The country where the flat is located."
				},
				"city": {
					"type": "string",
					"description": "The city or municipality where the flat is located."
				},
				"postalCode": {
					"type": "string",
					"description": "The postal or ZIP code for the flat's location."
				},
				"address": {
					"type": "string",
					"description": "The street address of the flat."
				},
				"roomCount": {
					"type": "number",
					"description": "The total number of rooms in the flat."
				},
				"squareMeters": {
					"type": "number",
					"description": "The total area of the flat in square meters."
				},
				"availabilityDate": {
					"type": "string",
					"format": "date",
					"description": "The date when the flat becomes available for occupancy."
				},
				"yearBuilt": {
					"type": "integer",
					"description": "The year when the building was constructed."
				},
				"hwgEnergyClass": {
					"type": ["string", "null"],
					"enum": ["A++", "A+", "A", "B", "C", "D", "E", "F", "G", null],
					"description": "The HWG energy efficiency classification of the flat."
				},
				"fgeeEnergyClass": {
					"type": ["string", "null"],
					"enum": ["A++", "A+", "A", "B", "C", "D", "E", "F", "G", null],
					"description": "The FGEE energy efficiency classification of the flat."
				},
				"rentPricePerMonth": {
					"type": ["number", "null"],
					"description": "The monthly rent price in euros. Required if the listing is for rent."
				},
					"cooperativeShare": {
					"type": ["number", "null"],
					"description": "The cooperative share or deposit required for renting, in euros. Applicable and required if the listing is for rent."
				},
					"salePrice": {
					"type": ["number", "null"],
					"description": "The sale price of the flat in euros. Required if the listing is for sale."
				},
					"additionalFees": {
					"type": ["number", "null"],
					"description": "Any additional fees associated with purchasing the flat, in euros. Optional and applicable if the listing is for sale."
				},
				"listingType": {
					"type": "string",
					"enum": ["rent", "sale", "both"],
					"description": "Indicates if the listing is available for rent, for sale, or both."
				},
				"detailsUrl": {
					"type": "string",
					"format": "uri",
					"description": "The URL to the detailed page for this flat listing on the housing cooperative's website."
				},
				"previewImageUrl": {
					"type": "string",
					"format": "uri",
					"description": "The URL to a preview image of the flat."
				}
				},
				"required": [
				"title",
				"housingCooperative",
				"projectId",
				"listingId",
				"country",
				"city",
				"postalCode",
				"address",
				"roomCount",
				"squareMeters",
				"availabilityDate",
				"yearBuilt",
				"listingType",
				"detailsUrl",
				"previewImageUrl"
				]
			},
			"description": "An array of flat listings."
			}
		},
		"required": [
			"scraperId",
			"timestamp",
			"listings"
		]
	}
	`
