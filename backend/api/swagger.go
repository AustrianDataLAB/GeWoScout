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
	response := models.NewInvokeResponse(http.StatusMovedPermanently, "Redirecting to Swagger UI")
	response.Outputs.Res.Headers["Location"] = "/api/swagger/index.html"
	render.JSON(w, r, response)
}

func SwaggerFileHandler(w http.ResponseWriter, r *http.Request) {
	req, err := models.InvokeRequestFromBody(r.Body)
	if err != nil {
		render.JSON(w, r, models.NewInvokeResponse(
			http.StatusBadRequest,
			models.Error{Message: err.Error(), StatusCode: http.StatusBadRequest},
		))
		return
	}

	r.Method = http.MethodGet
	r.URL, err = url.Parse(req.Data.Req.Url)
	if err != nil {
		render.JSON(w, r, models.NewInvokeResponse(
			http.StatusBadRequest,
			models.Error{Message: err.Error(), StatusCode: http.StatusBadRequest},
		))
		return
	}
	r.RequestURI = r.URL.Path

	mockWriter := NewMockResponseWriter()

	httpSwagger.Handler().ServeHTTP(mockWriter, r)

	//ir := models.NewInvokeResponse(mockWriter.StatusCode, mockWriter.Body.String())
	ir := models.InvokeResponse{}
	ir.Outputs.Res.StatusCode = mockWriter.StatusCode
	ir.Outputs.Res.Body = mockWriter.Body.String()
	ir.Outputs.Res.Headers = map[string]string{}
	for k, v := range mockWriter.Headers {
		ir.Outputs.Res.Headers[k] = strings.Join(v, "; ")
	}
	render.JSON(w, r, ir)
}
