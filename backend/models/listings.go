package models

type GetListingsResponse struct {
	Results           []Listing `json:"results"`
	ContinuationToken string    `json:"continuationToken"`
}
