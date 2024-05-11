package models

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

type Query struct {
	ContinuationToken *string `json:"continuationToken"`
	MinSize           *uint32 `json:"minSize,string" validation:"omitempty,gt=0"`
	MaxSize           *uint32 `json:"maxSize,string" validate:"omitempty,gtfieldcustom=MinSize,gt=0"`
	PageSize          *int    `json:"pageSize,string" validate:"omitempty,gt=0,lte=30"`
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
