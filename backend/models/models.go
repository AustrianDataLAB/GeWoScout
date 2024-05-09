package models

type Listing struct {
	ID                 string `json:"id"`
	PartitionKey       string `json:"_partitionKey"`
	Title              string `json:"title"`
	HousingCooperative string `json:"housingCooperative"`
	ProjectID          string `json:"projectId"`
	ListingID          string `json:"listingId"`
	Country            string `json:"country"`
	City               string `json:"city"`
	PostalCode         string `json:"postalCode"`
	Address            string `json:"address"`
	RoomCount          int    `json:"roomCount"`
	SquareMeters       int    `json:"squareMeters"`
	AvailabilityDate   string `json:"availabilityDate"`
	YearBuilt          int    `json:"yearBuilt"`
	HwgEnergyClass     string `json:"hwgEnergyClass"`
	FgeeEnergyClass    string `json:"fgeeEnergyClass"`
	ListingType        string `json:"listingType"`
	RentPricePerMonth  int    `json:"rentPricePerMonth"`
	CooperativeShare   int    `json:"cooperativeShare"`
	SalePrice          *int   `json:"salePrice"`
	AdditionalFees     *int   `json:"additionalFees"`
	DetailsURL         string `json:"detailsUrl"`
	PreviewImageURL    string `json:"previewImageUrl"`
	ScraperID          string `json:"scraperId"`
	CreatedAt          string `json:"createdAt"`
	LastModifiedAt     string `json:"lastModifiedAt"`
}

// Holds any form of error (either from Azure or some internal error)
type Error struct {
	Message string `json:"message"`
}
