package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/AustrianDataLAB/GeWoScout/backend/cosmos"
	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/go-chi/render"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func (h *Handler) HandleScraperResult(w http.ResponseWriter, r *http.Request) {
	injectedData := models.QueueBindingInput{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&injectedData); err != nil {
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusInternalServerError,
			nil,
			[]string{fmt.Sprintf("ScraperResultHandler | Failed to read invoke request body: %s", err.Error())},
		))
		return
	}

	msgId, err := strconv.Unquote(injectedData.Metadata.Id)
	if err != nil {
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusInternalServerError,
			nil,
			[]string{fmt.Sprintf("ScraperResultHandler %s | Failed to unquote message ID: %s", injectedData.Metadata.Id, err.Error())},
		))
		return
	}

	msgPlain, err := strconv.Unquote(injectedData.Data.Msg)
	if err != nil {
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusInternalServerError,
			nil,
			[]string{fmt.Sprintf("ScraperResultHandler %s | Failed to unquote message: %s", msgId, err.Error())},
		))
		return
	}

	logs := make([]string, 0)

	// TODO instead of returning 4XX, return 200 and send the message to a dead-letter queue
	// Validate the message
	msgLoader := gojsonschema.NewStringLoader(msgPlain)
	result, err := h.ScraperResultSchema.Validate(msgLoader)
	if err != nil {
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusUnprocessableEntity,
			nil,
			[]string{fmt.Sprintf("ScraperResultHandler %s | Failed to create message schema loader: %s", msgId, err.Error())},
		))
		return
	}

	if !result.Valid() {
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusInternalServerError,
			nil,
			[]string{fmt.Sprintf("ScraperResultHandler %s | Message validation failed: %s", msgId, result.Errors())},
		))
		return
	}

	scraperResult := models.ScraperResultList{}
	err = json.Unmarshal([]byte(msgPlain), &scraperResult)
	if err != nil {
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusUnprocessableEntity,
			nil,
			[]string{fmt.Sprintf("ScraperResultHandler %s | Failed to unmarshal message: %s", msgId, err.Error())},
		))
		return
	}

	idsByPk := make(map[string][]string)
	listingsById := make(map[string][]*models.Listing)
	for _, listing := range scraperResult.Listings {
		pk := mapPartitionKey(listing)
		idsByPk[pk] = append(idsByPk[pk], mapListingId(listing))
		listingsById[pk] = append(listingsById[pk], mapListing(listing, scraperResult.ScraperId, scraperResult.Timestamp))
	}

	logs = append(logs, fmt.Sprintf("ScraperResultHandler %s | Processing valid scraper result from %s with %d listings", msgId, scraperResult.ScraperId, len(scraperResult.Listings)))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nonExIds, err := cosmos.GetNonExistingIds(ctx, h.GetListingsByCityContainerClient(), idsByPk)
	if err != nil {
		render.JSON(w, r, models.NewHttpInvokeResponse(
			http.StatusInternalServerError,
			nil,
			[]string{fmt.Sprintf("ScraperResultHandler %s | Failed to get non-existing IDs: %s", msgId, err.Error())},
		))
		return
	}

	newCount := 0
	for _, v := range nonExIds {
		newCount += len(v)
	}

	if newCount == 0 {
		logs = append(logs, fmt.Sprintf("ScraperResultHandler %s | No new listings found", msgId))
	} else {
		logs = append(logs, fmt.Sprintf("ScraperResultHandler %s | Found %d new listings", msgId, newCount))
	}

	// Insert the new listings into the db
	// Path existing listings
	// Batch on partition key

	newListings := make([]*models.Listing, 0, newCount)
	container := h.GetListingsByCityContainerClient()

	for pk, ids := range idsByPk {
		listings := listingsById[pk]
		nonExIdsForPk := nonExIds[pk]

		partitionKey := azcosmos.NewPartitionKeyString(pk)
		batch := container.NewTransactionalBatch(partitionKey)

		for i, id := range ids {
			if slices.Contains(nonExIdsForPk, id) {
				marshalled, err := json.Marshal(listings[i])
				if err != nil {
					logs = append(logs, fmt.Sprintf("ScraperResultHandler %s | Failed to marshal listing %s: %s", msgId, id, err.Error()))
					continue
				}
				batch.UpsertItem(marshalled, nil)
				newListings = append(newListings, listings[i])
			} else {
				batch.PatchItem(id, createListingPatch(listings[i]), nil)
			}
		}

		// TODO do something with response??
		resp, err := container.ExecuteTransactionalBatch(ctx, batch, nil)
		if err != nil {
			logs = append(logs, fmt.Sprintf("ScraperResultHandler %s | Failed to execute batch for %s: %s", msgId, pk, err.Error()))
			//render.JSON(w, r, models.NewHttpInvokeResponse(http.StatusInternalServerError, struct{}{}))
			break // Go to end
		}

		if !resp.Success {
			// Transaction failed, look for the offending operation
			for index, operation := range resp.OperationResults {
				if operation.StatusCode != http.StatusFailedDependency {
					logs = append(logs, fmt.Sprintf("ScraperResultHandler %s | Transaction failed due to operation %v which failed with status code %d", msgId, index, operation.StatusCode))
					break
				}
			}
			break // Go to end
		}

		logs = append(logs, fmt.Sprintf("ScraperResultHandler %s | Inserted/patched batch of %d for %s", msgId, len(listings), pk))
	}

	newListingsOutput := make([]string, len(newListings))
	for i, l := range newListings {
		marshalled, _ := json.Marshal(l)
		newListingsOutput[i] = string(marshalled)
	}

	// For each non-existant ID, insert the listing and create a queue message
	invokeResponse := models.InvokeResponse{
		Logs: logs,
		Outputs: map[string]interface{}{
			"msgOut": newListingsOutput,
		},
	}
	render.JSON(w, r, invokeResponse)
}

func mapListing(scraperResult models.ScraperResultListing, scraperId string, timestamp time.Time) *models.Listing {
	return &models.Listing{
		ID:                 mapListingId(scraperResult),
		PartitionKey:       mapPartitionKey(scraperResult),
		Title:              scraperResult.Title,
		HousingCooperative: scraperResult.HousingCooperative,
		ProjectID:          scraperResult.ProjectId,
		ListingID:          scraperResult.ListingId,
		Country:            scraperResult.Country,
		City:               scraperResult.City,
		PostalCode:         scraperResult.PostalCode,
		Address:            scraperResult.Address,
		RoomCount:          scraperResult.RoomCount,
		SquareMeters:       scraperResult.SquareMeters,
		AvailabilityDate:   scraperResult.AvailabilityDate,
		YearBuilt:          scraperResult.YearBuilt,
		HwgEnergyClass:     scraperResult.HwgEnergyClass,
		FgeeEnergyClass:    scraperResult.FgeeEnergyClass,
		ListingType:        scraperResult.ListingType,
		RentPricePerMonth:  scraperResult.RentPricePerMonth,
		CooperativeShare:   scraperResult.CooperativeShare,
		SalePrice:          scraperResult.SalePrice,
		AdditionalFees:     scraperResult.AdditionalFees,
		DetailsURL:         scraperResult.DetailsUrl,
		PreviewImageURL:    scraperResult.PreviewImageUrl,
		ScraperID:          scraperId,
		CreatedAt:          timestamp,
		LastModifiedAt:     timestamp,
	}
}

func createListingPatch(l *models.Listing) azcosmos.PatchOperations {
	// TODO figure out which fields to actually patch
	patch := azcosmos.PatchOperations{}
	patch.AppendSet("/availabilityDate", l.AvailabilityDate)
	patch.AppendSet("/listingType", l.ListingType)
	patch.AppendSet("/rentPricePerMonth", l.RentPricePerMonth)
	patch.AppendSet("/cooperativeShare", l.CooperativeShare)
	patch.AppendSet("/salePrice", l.SalePrice)
	patch.AppendSet("/additionalFees", l.AdditionalFees)
	patch.AppendSet("/detailsUrl", l.DetailsURL)
	patch.AppendSet("/lastModifiedAt", l.LastModifiedAt)
	return patch
}

func mapListingId(listing models.ScraperResultListing) string {
	id := fmt.Sprintf("%s_%s_%s", listing.HousingCooperative, listing.ProjectId, listing.ListingId)
	id = strings.TrimSpace(id)
	id = removeDiacritics(id)
	id = strings.ToLower(id)
	id = strings.ReplaceAll(id, " ", "")
	return id
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
