package models

type GetListingsResponse struct {
	Results           []Listing `json:"results"`
	ContinuationToken *string   `json:"continuationToken"`
}

type ListingsQuery struct {
	Id                   string       `json:"id"`
	Email                string       `json:"email"`
	Title                *string      `json:"title" validate:"omitempty"`
	HousingCooperative   *string      `json:"housingCooperative" validate:"omitempty"`
	ProjectId            *string      `json:"projectId" validate:"omitempty"`
	PostalCode           *string      `json:"postalCode" validate:"omitempty"`
	RoomCount            *int         `json:"roomCount,string" validate:"omitempty,gt=0"`
	MinRoomCount         *int         `json:"minRoomCount,string" validate:"omitempty,gt=0"`
	MaxRoomCount         *int         `json:"maxRoomCount,string" validate:"omitempty,gt=0,gtfieldcustom=MinRoomCount"`
	MinSquareMeters      *int         `json:"minSqm,string" validate:"omitempty,gt=0"`
	MaxSquareMeters      *int         `json:"maxSqm,string" validate:"omitempty,gt=0,gtfieldcustom=MinSquareMeters"`
	AvailableFrom        *string      `json:"availableFrom" validate:"omitempty,datecustom"`
	MinYearBuilt         *int         `json:"minYearBuilt,string" validate:"omitempty,gt=1900"`
	MaxYearBuilt         *int         `json:"maxYearBuilt,string" validate:"omitempty,gt=1900,gtfieldcustom=MinYearBuilt"`
	MinHwgEnergyClass    *EnergyClass `json:"minHwgEnergyClass" validate:"omitempty,energycustom"`
	MinFgeeEnergyClass   *EnergyClass `json:"minFgeeEnergyClass" validate:"omitempty,energycustom"`
	ListingType          *ListingType `json:"listingType" validate:"omitempty,listingtypecustom"`
	MinRentPricePerMonth *int         `json:"minRentPrice,string" validate:"omitempty,gt=0"`
	MaxRentPricePerMonth *int         `json:"maxRentPrice,string" validate:"omitempty,gt=0,gtfieldcustom=MinRentPricePerMonth"`
	MinCooperativeShare  *int         `json:"minCooperativeShare,string" validate:"omitempty,gt=0"`
	MaxCooperativeShare  *int         `json:"maxCooperativeShare,string" validate:"omitempty,gt=0,gtfieldcustom=MinCooperativeShare"`
	MinSalePrice         *int         `json:"minSalePrice,string" validate:"omitempty,gt=0"`
	MaxSalePrice         *int         `json:"maxSalePrice,string" validate:"omitempty,gt=0,gtfieldcustom=MinSalePrice"`
	ContinuationToken    *string      `json:"continuationToken"`
	PageSize             *int         `json:"pageSize,string" validate:"omitempty,gt=0,lte=30"`
	SortBy               *string
	SortType             *SortType `json:"sortType" validate:"omitempty,sorttypecustom"`
}
