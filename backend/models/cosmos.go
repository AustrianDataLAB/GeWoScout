package models

type CosmosBindingInput struct {
	Data     Data     `json:"Data"`
	Metadata Metadata `json:"Metadata"`
}

type Data struct {
	Documents string   `json:"documents"`
	Metadata  Metadata `json:"Metadata"`
}

type Metadata struct {
	City string `json:"city"`
	ID   string `json:"id"`
}
