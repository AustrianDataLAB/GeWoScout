package main

import (
	"fmt"
	"log"
	"net/http"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!\n")
}

func getListings(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "List of Genossenschaftswohnungen\n")
}

func main() {
	port := "8080"
	log.Printf("About to listen on %s. Go to http://127.0.0.1:%s/", port, port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/listings", getListings)

	http.ListenAndServe(":"+port, mux)
}
