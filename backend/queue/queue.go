package queue

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"github.com/go-chi/render"
)

func QueueTriggerHandler(w http.ResponseWriter, r *http.Request) {
	var injectedData models.CosmosBindingInput

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&injectedData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: implement something
	md, _ := json.Marshal(injectedData.Metadata)
	d, _ := json.Marshal(injectedData.Data)
	log.Printf("Received metadata: %s\n", md)
	log.Printf("Received data: %s\n", d)

	invokeResponse := models.InvokeResponse{Logs: []string{}, ReturnValue: nil}
	render.JSON(w, r, invokeResponse)
}
