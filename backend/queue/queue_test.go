package queue

import (
	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"testing"
)

func TestRemoveDiacritics(t *testing.T) {
	input := "München über élégant Straße"
	expected := "Munchen uber elegant Strasse"

	actual := removeDiacritics(input)
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestMapPartitionKey(t *testing.T) {
	input := models.ScraperResultListing{
		City: "  St. Pölten ",
	}
	expected := "st.polten"

	actual := mapPartitionKey(input)
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
