package models

import (
	"encoding/json"
	"io"
	"time"

	"github.com/go-playground/validator/v10"
)

type ListingsQuery struct {
	Preferences

	ContinuationToken *string `json:"continuationToken"`
	PageSize          *int    `json:"pageSize,string" validate:"omitempty,gt=0,lte=30"`
	SortBy            *string
	SortType          *SortType `json:"sortType" validate:"omitempty,sorttypecustom"`
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

func dateCustomValidator(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().String())
	return err == nil
}

type InvokeRequest[Q any, B any] struct {
	Data struct {
		Req struct {
			Url        string
			Method     string
			Query      Q
			Headers    map[string]interface{}
			Host       []string
			UserAgent  []string `json:"User-Agent"`
			Params     map[string]string
			Identities []Identity
			Body       B
		} `json:"req"`
	}
	Metadata map[string]interface{}
}

type Identity struct {
	AuthenticationType *string `json:"AuthenticationType"`
	IsAuthenticated    bool    `json:"IsAuthenticated"`
	Actor              *string `json:"Actor"`
	BootstrapContext   *string `json:"BootstrapContext"`
	Label              *string `json:"Label"`
	Name               *string `json:"Name"`
	NameClaimType      string  `json:"NameClaimType"`
	RoleClaimType      string  `json:"RoleClaimType"`
}

func InvokeRequestFromBody[Q any, B any](body io.ReadCloser) (ir InvokeRequest[Q, B], err error) {
	var b []byte
	b, err = io.ReadAll(io.Reader(body))
	defer body.Close()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(b), &ir)
	if err != nil {
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("datecustom", dateCustomValidator)
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
