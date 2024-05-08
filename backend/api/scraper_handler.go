package api

import (
	"encoding/json"
	"fmt"
	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"github.com/go-chi/render"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

type bindingInput struct {
	Data     bindingData     `json:"Data"`
	Metadata bindingMetadata `json:"Metadata"`
}

type bindingData struct {
	Msg string `json:"msg"`
}

type bindingMetadata struct {
	Id           string `json:"Id"`
	DequeueCount string `json:"DequeueCount"`
}

func (h *Handler) CreateScraperResultHandler() http.HandlerFunc {

	schemaLoader := gojsonschema.NewStringLoader(models.ScraperResultListingSchema)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		log.Fatalf("Failed to create schema: %s", err.Error())
	}

	return func(w http.ResponseWriter, r *http.Request) {
		injectedData := bindingInput{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&injectedData); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, models.InvokeResponse{
				Logs: []string{fmt.Sprintf("ScraperResultHandler | Failed to read invoke request body: %s", err.Error())},
			})
			return
		}

		msgId, err := strconv.Unquote(injectedData.Metadata.Id)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, models.InvokeResponse{
				Logs: []string{fmt.Sprintf("ScraperResultHandler %s | Failed to unquote message ID: %s", injectedData.Metadata.Id, err.Error())},
			})
			return
		}

		msgPlain, err := strconv.Unquote(injectedData.Data.Msg)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, models.InvokeResponse{
				Logs: []string{fmt.Sprintf("ScraperResultHandler %s | Failed to unquote message: %s", msgId, err.Error())},
			})
			return
		}

		// TODO remove
		log.Printf("ScraperResultHandler %s | Received message: %s\n", msgId, msgPlain)

		// TODO instead of returning 4XX, return 200 and send the message to a dead-letter queue
		// Validate the message
		msgLoader := gojsonschema.NewStringLoader(msgPlain)
		result, err := schema.Validate(msgLoader)
		if err != nil {
			render.Status(r, http.StatusUnprocessableEntity)
			render.JSON(w, r, models.InvokeResponse{
				Logs: []string{fmt.Sprintf("ScraperResultHandler %s | Failed to create message schema loader: %s", msgId, err.Error())},
			})
			return
		}

		if !result.Valid() {
			render.JSON(w, r, models.InvokeResponse{
				Logs: []string{fmt.Sprintf("ScraperResultHandler %s | Message validation failed: %s", msgId, result.Errors())},
			})
			return
		}

		scraperResult := models.ScraperResultListing{}
		err = json.Unmarshal([]byte(msgPlain), &scraperResult)
		if err != nil {
			render.Status(r, http.StatusUnprocessableEntity)
			render.JSON(w, r, models.InvokeResponse{
				Logs: []string{fmt.Sprintf("ScraperResultHandler %s | Failed to unmarshal message: %s", msgId, err.Error())},
			})
			return
		}

		// TODO actually do something with the scraper results

		invokeResponse := models.InvokeResponse{Logs: []string{"TODO"}}
		render.JSON(w, r, invokeResponse)
	}
}

func mapListings(scraperResult models.ScraperResultListing) []models.Listing {
	// TODO
	return nil
}

func mapPartitionKey(listing models.ScraperResultListing) string {
	pk := listing.City
	pk = strings.TrimSpace(pk)
	pk = removeDiacritics(pk)
	pk = strings.ToLower(pk)
	pk = strings.ReplaceAll(pk, " ", "")
	return pk
}

// removeDiacritics removes diacritics from a string.
func removeDiacritics(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)

	// Additional replacement for special cases
	result = strings.ReplaceAll(result, "ß", "ss")

	return result
}
