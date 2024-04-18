package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AustrianDataLAB/GeWoScout/backend/api"
	"github.com/AustrianDataLAB/GeWoScout/backend/notification"
	"github.com/AustrianDataLAB/GeWoScout/backend/queue"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Alive"))
	})
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Alive"))
	})
	r.Get("/QueueTrigger", queue.QueueTriggerHandler)
	r.Get("/CosmosTrigger", notification.CosmosUpdateHandler)
	r.Get("/api/cities/{city}/listings", api.GetListings)
	// Mapping for /api/cities/{city}/listings/{id}
	// The Azure Function defined for this route has an injection from CosmosDB,
	// which means the original GET request is mapped to a POST request to this
	// route and the result is subsequently returned for the original GET
	// request.
	r.Post("/listingById", api.GetListingById)

	return r
}

func main() {
	port, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !exists {
		port = "8080"
	}

	log.Printf("About to listen on %s. Go to http://127.0.0.1:%s/", port, port)
	r := setupRouter()
	http.ListenAndServe(":"+port, r)
}
