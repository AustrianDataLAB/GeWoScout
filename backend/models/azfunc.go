package models

import "encoding/json"

type InvokeResponse struct {
	Outputs struct {
		Res struct {
			Body       string            `json:"body"`
			StatusCode int               `json:"statusCode"`
			Headers    map[string]string `json:"headers"`
		} `json:"res"`
	}
	Logs        []string
	ReturnValue interface{}
}

func NewInvokeResponse(statusCode int, body interface{}) (ir InvokeResponse) {
	d, _ := json.Marshal(body)
	ir.Outputs.Res.StatusCode = statusCode
	ir.Outputs.Res.Body = string(d)
	ir.Outputs.Res.Headers = map[string]string{
		"Content-Type": "application/json",
	}
	return
}
