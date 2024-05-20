package api

import (
	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"net/url"
	"strings"
)

func SwaggerBaseHandler(w http.ResponseWriter, r *http.Request) {
	response := models.NewHttpInvokeResponse(http.StatusMovedPermanently, "Redirecting to Swagger UI")
	httpOutput := response.Outputs["res"].(models.HttpResponse)
	httpOutput.Headers["Location"] = "/api/swagger/index.html"
	render.JSON(w, r, response)
}

func SwaggerFileHandler(w http.ResponseWriter, r *http.Request) {
	req, err := models.InvokeRequestFromBody(r.Body)
	if err != nil {
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusBadRequest,
			models.Error{Message: err.Error(), StatusCode: http.StatusBadRequest},
		))
		return
	}

	r.Method = http.MethodGet
	r.URL, err = url.Parse(req.Data.Req.Url)
	if err != nil {
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusBadRequest,
			models.Error{Message: err.Error(), StatusCode: http.StatusBadRequest},
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
