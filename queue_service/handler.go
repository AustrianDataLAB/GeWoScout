package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type InvokeResponse struct {
	Outputs     map[string]interface{}
	Logs        []string
	ReturnValue interface{}
}

type InvokeRequest struct {
	Data     map[string]interface{}
	Metadata map[string]interface{}
}

func queueTriggerHandler(w http.ResponseWriter, r *http.Request) {
	var invokeReq InvokeRequest

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&invokeReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: implement something
	md, _ := json.Marshal(invokeReq.Metadata)
	d, _ := json.Marshal(invokeReq.Data)
	log.Printf("Received metadata: %s\n", md)
	log.Printf("Received data: %s\n", d)

	invokeResponse := InvokeResponse{Logs: []string{}, ReturnValue: nil}

	resp, err := json.Marshal(invokeResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func main() {
	port, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if exists {
		fmt.Println("FUNCTIONS_CUSTOMHANDLER_PORT: " + port)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/QueueTrigger", queueTriggerHandler)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
