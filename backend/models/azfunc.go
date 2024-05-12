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
	RoomCount            *uint16      `json:"roomCount,string" validate:"omitempty"`
	MinRoomCount         *uint16      `json:"minRoomCount,string" validate:"omitempty"`
	MaxRoomCount         *uint16      `json:"maxRoomCount,string" validate:"omitempty,gtfieldcustom=MinRoomCount"`
	MinSquareMeters      *uint16      `json:"minSqm,string" validate:"omitempty"`
	MaxSquareMeters      *uint16      `json:"maxSqm,string" validate:"omitempty,gtfieldcustom=MinSquareMeters"`
	AvailableFrom        *time.Time   `json:"availableFrom,string" validate:"omitempty"`
	MinYearBuilt         *uint16      `json:"minYearBuilt,string" validate:"omitempty"`
	MaxYearBuilt         *uint16      `json:"maxYearBuilt,string" validate:"omitempty,gtfieldcustom=MinYearBuilt"`
	MinHwgEnergyClass    *EnergyClass `json:"minHwgEnergyClass" validate:"omitempty,energycustom"`
	MinFgeeEnergyClass   *EnergyClass `json:"minFgeeEnergyClass,string" validate:"omitempty,energycustom"`
	ListingType          *ListingType `json:"listingType" validate:"omitempty,listingtypecustom"`
	MinRentPricePerMonth *uint32      `json:"minRent,string" validate:"omitempty"`
	MaxRentPricePerMonth *uint32      `json:"maxRent,string" validate:"omitempty,gtfieldcustom=MinRentPricePerMonth"`
	MinCooperativeShare  *uint32      `json:"minCooperativeShare,string" validate:"omitempty"`
	MaxCooperativeShare  *uint32      `json:"maxCooperativeShare,string" validate:"omitempty,gtfieldcustom=MinCooperativeShare"`
	MinSalePrice         *uint32      `json:"minSalePrice,string" validate:"omitempty"`
	MaxSalePrice         *uint32      `json:"maxSalePrice,string" validate:"omitempty,gtfieldcustom=MinSalePrice"`
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
		return otherField.Elem().Uint() <= fl.Field().Uint()
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
