package models

import (
	"encoding/json"
	"io"
)

type InvokeRequest struct {
	Data struct {
		Req struct {
			Url        string
			Method     string
			Query      map[string]string
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
	return
}

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
