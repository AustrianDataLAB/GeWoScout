package notification

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AustrianDataLAB/GeWoScout/backend/models"
)

func CosmosUpdateHandler(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	fmt.Println(t.Month())
	fmt.Println(t.Day())
	fmt.Println(t.Year())
	fmt.Println(r.Header)
	ua := r.Header.Get("User-Agent")
	fmt.Printf("user agent is: %s \n", ua)
	invocationid := r.Header.Get("X-Azure-Functions-InvocationId")
	fmt.Printf("invocationid is: %s \n", invocationid)

	queryParams := r.URL.Query()

	for k, v := range queryParams {
		fmt.Println("k:", k, "v:", v)
	}

	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	fmt.Println("Body:", string(body))

	var cosmosTrigger models.CosmosBindingInput
	err := json.Unmarshal(body, &cosmosTrigger)
	if err != nil {
		fmt.Println("Error unmarshalling JSON trigger: ", err)
	}

	fmt.Println("CosmosTrigger:", cosmosTrigger)

	md, _ := json.Marshal(cosmosTrigger.Metadata)
	d, _ := json.Marshal(cosmosTrigger.Data)

	log.Printf("Metadata: %s\n", md)
	log.Printf("Data: %s\n", d)

	var dataStr string
	err = json.Unmarshal([]byte(cosmosTrigger.Data.Documents), &dataStr)
	if err != nil {
		fmt.Println("Error unmarshalling JSON string: ", err)
	}

	// `MyCosmosDocument` contains an escaped JSON string, parse it too.
	var documents []map[string]interface{} // Using a map for flexible structure handling
	err = json.Unmarshal([]byte(dataStr), &documents)
	if err != nil {
		fmt.Println("Error unmarshalling JSON documents: ", err)
	}
	fmt.Printf("Parsed Documents: %+v\n", documents)

	invokeResponse := models.InvokeResponse{Logs: []string{}, ReturnValue: documents[0]["id"].(string)}

	js, err := json.Marshal(invokeResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
