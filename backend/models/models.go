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
	RentPricePerMonth  *int   `json:"rentPricePerMonth"`
	CooperativeShare   *int   `json:"cooperativeShare"`
	SalePrice          *int   `json:"salePrice"`
	AdditionalFees     *int   `json:"additionalFees"`
	DetailsURL         string `json:"detailsUrl"`
	PreviewImageURL    string `json:"previewImageUrl"`
	ScraperID          string `json:"scraperId"`
	CreatedAt          string `json:"createdAt"`
	LastModifiedAt     string `json:"lastModifiedAt"`
}

type EnergyClass string

func (c EnergyClass) IsEnumValue() bool {
	return (c == EnergyClassAplusplus ||
		c == EnergyClassAplus ||
		c == EnergyClassA ||
		c == EnergyClassB ||
		c == EnergyClassC ||
		c == EnergyClassD ||
		c == EnergyClassE ||
		c == EnergyClassF)
}

type ListingType string

func (t ListingType) IsEnumValue() bool {
	return t == ListingTypeRent || t == ListingTypeSale || t == ListingTypeBoth
}

type SortType string

func (t SortType) IsEnumValue() bool {
	return t == SortTypeAsc || t == SortTypeDesc
}

const (
	EnergyClassAplusplus EnergyClass = "A++"
	EnergyClassAplus     EnergyClass = "A+"
	EnergyClassA         EnergyClass = "A"
	EnergyClassB         EnergyClass = "B"
	EnergyClassC         EnergyClass = "C"
	EnergyClassD         EnergyClass = "D"
	EnergyClassE         EnergyClass = "E"
	EnergyClassF         EnergyClass = "F"

	ListingTypeRent ListingType = "rent"
	ListingTypeSale ListingType = "sale"
	ListingTypeBoth ListingType = "both"

	SortTypeAsc  SortType = "ASC"
	SortTypeDesc SortType = "DESC"
)

type StringEnum interface {
	IsEnumValue() bool
}

// Holds any form of error (either from Azure or some internal error)
type Error struct {
	Message string `json:"message"`
}
