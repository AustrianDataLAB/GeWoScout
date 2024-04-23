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
	r := setupRouter()

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
