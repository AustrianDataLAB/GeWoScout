package notification

import (
	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"testing"
	"time"
)

func TestGenerateEmailContentWithMinimalListing(t *testing.T) {
	listing := models.Listing{}

	emailContent, err := generateEmailContent(listing)

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if emailContent == "" {
		t.Error("Expected email content, got empty string")
	}
}

func TestGenerateEmailContentWithFullListing(t *testing.T) {
	listing := models.Listing{
		ID:                 "123",
		PartitionKey:       "123",
		Title:              "Test Title",
		HousingCooperative: "Test Housing Cooperative",
		ProjectID:          "123",
		ListingID:          "123",
		Country:            "Test Country",
		City:               "Test City",
		PostalCode:         "123",
		Address:            "Test Address",
		RoomCount:          1,
		SquareMeters:       1,
		AvailabilityDate:   "Test Availability Date",
		YearBuilt:          new(int),
		HwgEnergyClass:     new(string),
		FgeeEnergyClass:    new(string),
		ListingType:        "Test Listing Type",
		RentPricePerMonth:  new(int),
		CooperativeShare:   new(int),
		SalePrice:          new(int),
		AdditionalFees:     new(int),
		DetailsURL:         "Test Details URL",
		PreviewImageURL:    "Test Preview Image URL",
		ScraperID:          "123",
		CreatedAt:          time.Now(),
		LastModifiedAt:     time.Now(),
	}

	emailContent, err := generateEmailContent(listing)

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if emailContent == "" {
		t.Error("Expected email content, got empty string")
	}
}
