package models

import (
	"encoding/json"
)

type InvokeRequest[Q any] struct {
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
			Body       string
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
