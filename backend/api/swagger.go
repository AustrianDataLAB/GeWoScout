package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SwaggerBaseHandler(w http.ResponseWriter, r *http.Request) {
	response := models.NewHttpInvokeResponse(http.StatusMovedPermanently, "Redirecting to Swagger UI", nil)
	httpOutput := response.Outputs["res"].(models.HttpResponse)
	httpOutput.Headers["Location"] = "/api/swagger/index.html"
	render.JSON(w, r, response)
}

func SwaggerFileHandler(w http.ResponseWriter, r *http.Request) {
	req, err := models.InvokeRequestFromBody[interface{}, interface{}](r.Body)
	if err != nil {
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusBadRequest,
			models.Error{Message: err.Error()},
			[]string{err.Error()},
		))
		return
	}

	r.Method = http.MethodGet
	r.URL, err = url.Parse(req.Data.Req.Url)
	if err != nil {
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusBadRequest,
			models.Error{Message: err.Error()},
			[]string{err.Error()},
		))
		return
	}
	r.RequestURI = r.URL.Path

	mockWriter := NewMockResponseWriter()

	httpSwagger.Handler().ServeHTTP(mockWriter, r)

	//ir := models.NewHttpInvokeResponse(mockWriter.StatusCode, mockWriter.Body.String())
	ir := models.InvokeResponse{}
	ir.Outputs = map[string]interface{}{}
	headers := make(map[string]string)
	for k, v := range mockWriter.Headers {
		headers[k] = strings.Join(v, "; ")
	}
	ir.Outputs["res"] = models.HttpResponse{
		StatusCode: mockWriter.StatusCode,
		Body:       mockWriter.Body.String(),
		Headers:    headers,
	}
	render.JSON(w, r, ir)
}
