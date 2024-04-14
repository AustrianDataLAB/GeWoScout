package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type ReturnValue struct {
	Data string
}
type InvokeResponse struct {
	Outputs     map[string]interface{}
	Logs        []string
	ReturnValue interface{}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
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

	w.Write([]byte("Hello World from go worker"))
}

type CosmosData struct {
	Documents string `json:"documents"`
}

type CosmosMetadataSys struct {
	MethodName string `json:"MethodName"`
	UtcNow     string `json:"UtcNow"`
	RandGuid   string `json:"RandGuid"`
}

type CosmosMetadata struct {
	Sys CosmosMetadataSys `json:"sys"`
}

type CosmosTrigger struct {
	Data     CosmosData     `json:"Data"`
	Metadata CosmosMetadata `json:"Metadata"`
}

func cosmosUpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	var cosmosTrigger CosmosTrigger
	err := json.Unmarshal(body, &cosmosTrigger)
	if err != nil {
		fmt.Println("Error unmarshalling JSON trigger: ", err)
	}

	fmt.Println("CosmosTrigger:", cosmosTrigger)

	fmt.Println("Cosmos Trigger:")
	fmt.Println("Metadata", cosmosTrigger.Metadata.Sys.MethodName, cosmosTrigger.Metadata.Sys.UtcNow, cosmosTrigger.Metadata.Sys.RandGuid)
	fmt.Println("Data", cosmosTrigger.Data.Documents)

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

	returnValue := ReturnValue{Data: "Hello from Go!"}
	invokeResponse := InvokeResponse{Logs: []string{"test log1", "test log2"}, ReturnValue: returnValue}

	js, err := json.Marshal(invokeResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	customHandlerPort, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if exists {
		fmt.Println("FUNCTIONS_CUSTOMHANDLER_PORT: " + customHandlerPort)
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/api/hello", helloHandler)
	mux.HandleFunc("/CosmosTrigger", cosmosUpdateHandler)

	fmt.Println("Go server Listening...on FUNCTIONS_CUSTOMHANDLER_PORT:", customHandlerPort)
	log.Fatal(http.ListenAndServe(":"+customHandlerPort, mux))
}
