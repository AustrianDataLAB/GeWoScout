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

func GetEnergyClasses() []EnergyClass {
	return []EnergyClass{
		EnergyClassAplusplus,
		EnergyClassAplus,
		EnergyClassA,
		EnergyClassB,
		EnergyClassC,
		EnergyClassD,
		EnergyClassE,
		EnergyClassF,
	}
}

func (c EnergyClass) IsEnumValue() bool {
	return c.GetIndex() != -1
}

func (c EnergyClass) GetIndex() int {
	arr := GetEnergyClasses()
	for i := 0; i < len(arr); i++ {
		if arr[i] == c {
			return i
		}
	}
	return -1
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

type NotificationSettings struct {
	Preferences
}

type Preferences struct {
	PartitionKey         string       `json:"_partitionKey,omitempty"`
	Id                   string       `json:"id"`
	Email                string       `json:"email"`
	Title                *string      `json:"title" validate:"omitempty"`
	HousingCooperative   *string      `json:"housingCooperative" validate:"omitempty"`
	ProjectId            *string      `json:"projectId" validate:"omitempty"`
	PostalCode           *string      `json:"postalCode" validate:"omitempty"`
	RoomCount            *int         `json:"roomCount" validate:"omitempty,gt=0"`
	MinRoomCount         *int         `json:"minRoomCount" validate:"omitempty,gt=0"`
	MaxRoomCount         *int         `json:"maxRoomCount" validate:"omitempty,gt=0,gtfieldcustom=MinRoomCount"`
	MinSquareMeters      *int         `json:"minSqm" validate:"omitempty,gt=0"`
	MaxSquareMeters      *int         `json:"maxSqm" validate:"omitempty,gt=0,gtfieldcustom=MinSquareMeters"`
	AvailableFrom        *string      `json:"availableFrom" validate:"omitempty,datecustom"`
	MinYearBuilt         *int         `json:"minYearBuilt" validate:"omitempty,gt=1900"`
	MaxYearBuilt         *int         `json:"maxYearBuilt" validate:"omitempty,gt=1900,gtfieldcustom=MinYearBuilt"`
	MinHwgEnergyClass    *EnergyClass `json:"minHwgEnergyClass" validate:"omitempty,energycustom"`
	MinFgeeEnergyClass   *EnergyClass `json:"minFgeeEnergyClass" validate:"omitempty,energycustom"`
	ListingType          *ListingType `json:"listingType" validate:"omitempty,listingtypecustom"`
	MinRentPricePerMonth *int         `json:"minRentPrice" validate:"omitempty,gt=0"`
	MaxRentPricePerMonth *int         `json:"maxRentPrice" validate:"omitempty,gt=0,gtfieldcustom=MinRentPricePerMonth"`
	MinCooperativeShare  *int         `json:"minCooperativeShare" validate:"omitempty,gt=0"`
	MaxCooperativeShare  *int         `json:"maxCooperativeShare" validate:"omitempty,gt=0,gtfieldcustom=MinCooperativeShare"`
	MinSalePrice         *int         `json:"minSalePrice" validate:"omitempty,gt=0"`
	MaxSalePrice         *int         `json:"maxSalePrice" validate:"omitempty,gt=0,gtfieldcustom=MinSalePrice"`
}
