package models

import (
	"encoding/json"
	"io"
	"time"

	"github.com/go-playground/validator/v10"
)

type Query struct {
	ContinuationToken    *string      `json:"continuationToken"`
	PageSize             *int         `json:"pageSize,string" validate:"omitempty,gt=0,lte=30"`
	Title                *string      `json:"title,string" validate:"omitempty"`
	HousingCooperative   *string      `json:"housingCooperative,string" validate:"omitempty"`
	ProjectId            *string      `json:"projectId,string" validate:"omitempty"`
	PostalCode           *string      `json:"postalCode" validate:"omitempty"`
	RoomCount            *int         `json:"roomCount,string" validate:"omitempty,gt=0"`
	MinRoomCount         *int         `json:"minRoomCount,string" validate:"omitempty,gt=0"`
	MaxRoomCount         *int         `json:"maxRoomCount,string" validate:"omitempty,gt=0,gtfieldcustom=MinRoomCount"`
	MinSquareMeters      *int         `json:"minSqm,string" validate:"omitempty,gt=0"`
	MaxSquareMeters      *int         `json:"maxSqm,string" validate:"omitempty,gt=0,gtfieldcustom=MinSquareMeters"`
	AvailableFrom        *time.Time   `json:"availableFrom" validate:"omitempty"`
	MinYearBuilt         *int         `json:"minYearBuilt,string" validate:"omitempty,gt=1900"`
	MaxYearBuilt         *int         `json:"maxYearBuilt,string" validate:"omitempty,gt=1900,gtfieldcustom=MinYearBuilt"`
	MinHwgEnergyClass    *EnergyClass `json:"minHwgEnergyClass" validate:"omitempty,energycustom"`
	MinFgeeEnergyClass   *EnergyClass `json:"minFgeeEnergyClass,string" validate:"omitempty,energycustom"`
	ListingType          *ListingType `json:"listingType" validate:"omitempty,listingtypecustom"`
	MinRentPricePerMonth *int         `json:"minRent,string" validate:"omitempty,gt=0"`
	MaxRentPricePerMonth *int         `json:"maxRent,string" validate:"omitempty,gt=0,gtfieldcustom=MinRentPricePerMonth"`
	MinCooperativeShare  *int         `json:"minCooperativeShare,string" validate:"omitempty,gt=0"`
	MaxCooperativeShare  *int         `json:"maxCooperativeShare,string" validate:"omitempty,gt=0,gtfieldcustom=MinCooperativeShare"`
	MinSalePrice         *int         `json:"minSalePrice,string" validate:"omitempty,gt=0"`
	MaxSalePrice         *int         `json:"maxSalePrice,string" validate:"omitempty,gt=0,gtfieldcustom=MinSalePrice"`
	SortBy               *string
	SortType             *SortType `json:"sortType" validate:"omitempty,sorttypecustom"`
}

func enumFieldValidator[T StringEnum](fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(T)
	return value.IsEnumValue()
}

func gtFieldIgnoreNilValidator(fl validator.FieldLevel) bool {
	otherField := fl.Parent().FieldByName(fl.Param())
	if !otherField.IsNil() {
		return otherField.Elem().Int() <= fl.Field().Int()
	}
	return true
}

type InvokeRequest struct {
	Data struct {
		Req struct {
			Url        string
			Method     string
			Query      Query
			Headers    map[string]interface{}
			Host       []string
			UserAgent  []string `json:"User-Agent"`
			Params     map[string]string
			Identities []interface{}
		} `json:"req"`
	}
	Metadata map[string]interface{}
}

func InvokeRequestFromBody(body io.ReadCloser) (ir InvokeRequest, err error) {
	var b []byte
	b, err = io.ReadAll(io.Reader(body))
	defer body.Close()
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &ir)
	if err != nil {
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("gtfieldcustom", gtFieldIgnoreNilValidator)
	validate.RegisterValidation("energycustom", enumFieldValidator[EnergyClass])
	validate.RegisterValidation("listingtypecustom", enumFieldValidator[ListingType])
	validate.RegisterValidation("sorttypecustom", enumFieldValidator[SortType])
	err = validate.Struct(ir)
	return
}

type HttpResponse struct {
	Body       string            `json:"body"`
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
}

type InvokeResponse struct {
	Outputs     map[string]interface{}
	Logs        []string
	ReturnValue interface{}
}

func NewHttpInvokeResponse(statusCode int, body interface{}, logs []string) (ir InvokeResponse) {
	d, _ := json.Marshal(body)
	ir.Outputs = make(map[string]interface{})
	ir.Outputs["res"] = HttpResponse{
		StatusCode: statusCode,
		Body:       string(d),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	ir.Logs = logs
	return
}

func NewHttpInvokeResponseWithHeaders(statusCode int, body interface{}, headers map[string]string, logs []string) (ir InvokeResponse) {
	d, _ := json.Marshal(body)
	ir.Outputs = make(map[string]interface{})
	ir.Outputs["res"] = HttpResponse{
		StatusCode: statusCode,
		Body:       string(d),
		Headers:    headers,
	}
	ir.Logs = logs
	return
}
