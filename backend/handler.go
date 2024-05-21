package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AustrianDataLAB/GeWoScout/backend/api"
	_ "github.com/AustrianDataLAB/GeWoScout/backend/docs"
	"github.com/AustrianDataLAB/GeWoScout/backend/notification"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func setupRouter(useSwagger bool) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	h := api.NewHandler()

	r.Get("/", h.HandleHealth)
	r.Post("/health", h.HandleHealth)
	r.Post("/scraperResultTrigger", h.HandleScraperResult)
	r.Post("/CosmosTrigger", notification.CosmosUpdateHandler)
	r.Post("/listings", h.GetListings)
	// Mapping for /api/cities/{city}/listings/{id}
	// The Azure Function defined for this route has an injection from CosmosDB,
	// which means the original GET request is mapped to a POST request to this
	// route and the result is subsequently returned for the original GET
	// request.
	r.Post("/listingById", h.GetListingById)

	if useSwagger {
		r.Post("/swagger", api.SwaggerBaseHandler)
		r.Post("/swaggerFiles", api.SwaggerFileHandler)
	}

	return r
}

// @title GeWoScout API
// @version 1
// @description This is the API for the GeWoScout project.
// @BasePath /api
func main() {
	port, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !exists {
		port = "8080"
	}

	log.Printf("About to listen on %s. Go to http://127.0.0.1:%s/", port, port)

	// TODO don't use Swagger in production??
	r := setupRouter(true)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
