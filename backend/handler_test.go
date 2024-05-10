package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AustrianDataLAB/GeWoScout/backend/models"
)

func Test(t *testing.T) {
	r := setupRouter(false)

	bodyStr := `{"Data":{"documents":"\"[{\\\"id\\\":\\\"ArtHabitat_AHProj055_AHFlat681\\\",\\\"_partitionKey\\\":\\\"salzburg\\\",\\\"title\\\":\\\"Artistic 2-Bedroom Close to Historic Center\\\",\\\"housingCooperative\\\":\\\"Art Habitat\\\",\\\"projectId\\\":\\\"AHProj055\\\",\\\"listingId\\\":\\\"AHFlat681\\\",\\\"country\\\":\\\"Austria\\\",\\\"city\\\":\\\"Salzburg\\\",\\\"postalCode\\\":\\\"5020\\\",\\\"address\\\":\\\"Getreidegasse 47\\\",\\\"roomCount\\\":2,\\\"squareMeters\\\":75,\\\"availabilityDate\\\":\\\"2024-07-15\\\",\\\"yearBuilt\\\":1985,\\\"hwgEnergyClass\\\":\\\"C\\\",\\\"fgeeEnergyClass\\\":\\\"B\\\",\\\"listingType\\\":\\\"rent\\\",\\\"rentPricePerMonth\\\":950,\\\"cooperativeShare\\\":3000,\\\"salePrice\\\":null,\\\"additionalFees\\\":null,\\\"detailsUrl\\\":\\\"https://arthabitat.at/listings/AHFlat681\\\",\\\"previewImageUrl\\\":\\\"https://arthabitat.at/listings/AHFlat681/preview.jpg\\\",\\\"scraperId\\\":\\\"salzburgHousingScraper002\\\",\\\"createdAt\\\":\\\"2024-04-14T14:30:00Z\\\",\\\"lastModifiedAt\\\":\\\"2024-04-14T14:30:00Z\\\",\\\"_rid\\\":\\\"9TZrALpFOV8KAAAAAAAAAA==\\\",\\\"_self\\\":\\\"dbs/9TZrAA==/colls/9TZrALpFOV8=/docs/9TZrALpFOV8KAAAAAAAAAA==/\\\",\\\"_etag\\\":\\\"\\\\\\\"0000fe06-0000-1500-0000-661c2c970000\\\\\\\"\\\",\\\"_attachments\\\":\\\"attachments/\\\",\\\"_ts\\\":1713122455,\\\"_lsn\\\":94}]\""},"Metadata":{"sys":{"MethodName":"CosmosTrigger","UtcNow":"2024-04-14T19:21:00.4080539Z","RandGuid":"da961a13-c6a0-4e75-b57b-367b13353311"}}}`

	req, _ := http.NewRequest("POST", "/CosmosTrigger", strings.NewReader(bodyStr))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type to be application/json")
	}

	// Note: Doesn't really make sense to check here, but would make sense if there is actually e.g. a http response
	response := models.InvokeResponse{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("error unmarshalling response: %v", err)
		return
	}

	if response.ReturnValue == nil {
		t.Errorf("expected ReturnValue to be non-nil")
		return
	}

	returnVal, ok := response.ReturnValue.(string)
	if !ok {
		t.Errorf("expected ReturnValue to be of type string")
		return
	}

	if returnVal != "ArtHabitat_AHProj055_AHFlat681" {
		t.Errorf("expected returnVal to be ArtHabitat_AHProj055_AHFlat681, got %s", returnVal)
	}
}

func TestGetListingById(t *testing.T) {

}

func TestHandleScraperResult(t *testing.T) {
	r := setupRouter(false)

	bodyStr := `{"Data":{"msg":"\"{\\\"scraperId\\\":\\\"viennaHousingScraper002\\\",\\\"timestamp\\\":\\\"2024-04-06T15:30:00Z\\\",\\\"listings\\\":[{\\\"title\\\":\\\"Modern3-BedroomApartmentinCentralVienna\\\",\\\"housingCooperative\\\":\\\"FutureLivingGenossenschaft\\\",\\\"projectId\\\":\\\"FLG2024\\\",\\\"listingId\\\":\\\"12345ABC\\\",\\\"country\\\":\\\"Austria\\\",\\\"city\\\":\\\"Vienna\\\",\\\"postalCode\\\":\\\"1010\\\",\\\"address\\\":\\\"Beispielgasse42\\\",\\\"roomCount\\\":3,\\\"squareMeters\\\":95,\\\"availabilityDate\\\":\\\"2024-09-01\\\",\\\"yearBuilt\\\":2019,\\\"hwgEnergyClass\\\":\\\"A\\\",\\\"fgeeEnergyClass\\\":\\\"A+\\\",\\\"listingType\\\":\\\"both\\\",\\\"rentPricePerMonth\\\":1200,\\\"cooperativeShare\\\":5000,\\\"salePrice\\\":350000,\\\"additionalFees\\\":6500,\\\"detailsUrl\\\":\\\"https://www.futurelivinggenossenschaft.at/listings/12345ABC\\\",\\\"previewImageUrl\\\":\\\"https://www.futurelivinggenossenschaft.at/listings/12345ABC/preview.jpg\\\"}]}\""},"Metadata":{"DequeueCount":"5","ExpirationTime":"2024-05-17T16:13:45+00:00","Id":"\"4434c659-815f-434e-a851-f1df7b701f27\"","InsertionTime":"2024-05-10T16:13:45+00:00","NextVisibleTime":"2024-05-10T16:23:46+00:00","PopReceipt":"\"AgAAAAMAAAAAAAAAVqQNb/ai2gE=\"","sys":{"MethodName":"scraperResultTrigger","UtcNow":"2024-05-10T16:13:47.1756352Z","RandGuid":"828f9d97-f30d-46ff-a609-024f97f376d8"},"scraperId":"\"viennaHousingScraper002\"","timestamp":"\"04/06/2024 15:30:00\""}}`

	req, _ := http.NewRequest("POST", "/scraperResultTrigger", strings.NewReader(bodyStr))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Errorf("expected Content-Type to be application/json; charset=utf-8")
	}

}
