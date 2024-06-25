package models

type Listing struct {
	ID                 string   `json:"id"`
	PartitionKey       string   `json:"_partitionKey"`
	Title              string   `json:"title"`
	HousingCooperative string   `json:"housingCooperative"`
	ProjectID          string   `json:"projectId"`
	ListingID          string   `json:"listingId"`
	Country            string   `json:"country"`
	City               string   `json:"city"`
	PostalCode         string   `json:"postalCode"`
	Address            string   `json:"address"`
	RoomCount          *float32 `json:"roomCount"`
	SquareMeters       float32  `json:"squareMeters"`
	AvailabilityDate   string   `json:"availabilityDate"`
	YearBuilt          *int     `json:"yearBuilt"`
	HwgEnergyClass     *string  `json:"hwgEnergyClass"`
	FgeeEnergyClass    *string  `json:"fgeeEnergyClass"`
	ListingType        string   `json:"listingType"`
	RentPricePerMonth  *float32 `json:"rentPricePerMonth"`
	CooperativeShare   *float32 `json:"cooperativeShare"`
	SalePrice          *float32 `json:"salePrice"`
	AdditionalFees     *float32 `json:"additionalFees"`
	DetailsURL         string   `json:"detailsUrl"`
	PreviewImageURL    string   `json:"previewImageUrl"`
	ScraperID          string   `json:"scraperId"`
	CreatedAt          string   `json:"createdAt"`
	LastModifiedAt     string   `json:"lastModifiedAt"`
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
		EnergyClassG,
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
	EnergyClassG         EnergyClass = "G"

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
	City *string `json:"city"`
	Preferences
}

type UserData struct {
	City *string `json:"city"`
	Preferences
}

type Preferences struct {
	PartitionKey         string       `json:"_partitionKey,omitempty"`
	Id                   string       `json:"id"`
	Email                string       `json:"email" validate:"email"`
	Title                *string      `json:"title,omitempty" validate:"omitempty"`
	HousingCooperative   *string      `json:"housingCooperative,omitempty" validate:"omitempty"`
	ProjectId            *string      `json:"projectId,omitempty" validate:"omitempty"`
	PostalCode           *string      `json:"postalCode,omitempty" validate:"omitempty"`
	RoomCount            *float32     `json:"roomCount,omitempty" validate:"omitempty,gt=0"`
	MinRoomCount         *float32     `json:"minRoomCount,omitempty" validate:"omitempty,gt=0"`
	MaxRoomCount         *float32     `json:"maxRoomCount,omitempty" validate:"omitempty,gt=0,gtfieldcustom=MinRoomCount"`
	MinSquareMeters      *float32     `json:"minSqm,omitempty" validate:"omitempty,gt=0"`
	MaxSquareMeters      *float32     `json:"maxSqm,omitempty" validate:"omitempty,gt=0,gtfieldcustom=MinSquareMeters"`
	AvailableFrom        *string      `json:"availableFrom,omitempty" validate:"omitempty,datecustom"`
	MinYearBuilt         *int         `json:"minYearBuilt,omitempty" validate:"omitempty,gt=1900"`
	MaxYearBuilt         *int         `json:"maxYearBuilt,omitempty" validate:"omitempty,gt=1900,gtfieldcustom=MinYearBuilt"`
	MinHwgEnergyClass    *EnergyClass `json:"minHwgEnergyClass,omitempty" validate:"omitempty,energycustom"`
	MinFgeeEnergyClass   *EnergyClass `json:"minFgeeEnergyClass,omitempty" validate:"omitempty,energycustom"`
	ListingType          *ListingType `json:"listingType,omitempty" validate:"omitempty,listingtypecustom"`
	MinRentPricePerMonth *float32     `json:"minRentPrice,omitempty" validate:"omitempty,gt=0"`
	MaxRentPricePerMonth *float32     `json:"maxRentPrice,omitempty" validate:"omitempty,gt=0,gtfieldcustom=MinRentPricePerMonth"`
	MinCooperativeShare  *float32     `json:"minCooperativeShare,omitempty" validate:"omitempty,gt=0"`
	MaxCooperativeShare  *float32     `json:"maxCooperativeShare,omitempty" validate:"omitempty,gt=0,gtfieldcustom=MinCooperativeShare"`
	MinSalePrice         *float32     `json:"minSalePrice,omitempty" validate:"omitempty,gt=0"`
	MaxSalePrice         *float32     `json:"maxSalePrice,omitempty" validate:"omitempty,gt=0,gtfieldcustom=MinSalePrice"`
}

type NotificationData struct {
	Emails  []string `json:"emails"`
	Listing Listing  `json:"listing"`
}
