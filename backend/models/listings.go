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
	RoomCount            *float32     `json:"roomCount,string" validate:"omitempty,gt=0"`
	MinRoomCount         *float32     `json:"minRoomCount,string" validate:"omitempty,gt=0"`
	MaxRoomCount         *float32     `json:"maxRoomCount,string" validate:"omitempty,gt=0,gtfieldcustom=MinRoomCount"`
	MinSquareMeters      *float32     `json:"minSqm,string" validate:"omitempty,gt=0"`
	MaxSquareMeters      *float32     `json:"maxSqm,string" validate:"omitempty,gt=0,gtfieldcustom=MinSquareMeters"`
	AvailableFrom        *string      `json:"availableFrom" validate:"omitempty,datecustom"`
	MinYearBuilt         *int         `json:"minYearBuilt,string" validate:"omitempty,gt=1900"`
	MaxYearBuilt         *int         `json:"maxYearBuilt,string" validate:"omitempty,gt=1900,gtfieldcustom=MinYearBuilt"`
	MinHwgEnergyClass    *EnergyClass `json:"minHwgEnergyClass" validate:"omitempty,energycustom"`
	MinFgeeEnergyClass   *EnergyClass `json:"minFgeeEnergyClass" validate:"omitempty,energycustom"`
	ListingType          *ListingType `json:"listingType" validate:"omitempty,listingtypecustom"`
	MinRentPricePerMonth *float32     `json:"minRentPrice,string" validate:"omitempty,gt=0"`
	MaxRentPricePerMonth *float32     `json:"maxRentPrice,string" validate:"omitempty,gt=0,gtfieldcustom=MinRentPricePerMonth"`
	MinCooperativeShare  *float32     `json:"minCooperativeShare,string" validate:"omitempty,gt=0"`
	MaxCooperativeShare  *float32     `json:"maxCooperativeShare,string" validate:"omitempty,gt=0,gtfieldcustom=MinCooperativeShare"`
	MinSalePrice         *float32     `json:"minSalePrice,string" validate:"omitempty,gt=0"`
	MaxSalePrice         *float32     `json:"maxSalePrice,string" validate:"omitempty,gt=0,gtfieldcustom=MinSalePrice"`
	ContinuationToken    *string      `json:"continuationToken"`
	PageSize             *int         `json:"pageSize,string" validate:"omitempty,gt=0,lte=30"`
	SortBy               *string
	SortType             *SortType `json:"sortType" validate:"omitempty,sorttypecustom"`
}
