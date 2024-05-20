package models

type QueueBindingInput struct {
	Data     QueueBindingData     `json:"Data"`
	Metadata QueueBindingMetadata `json:"Metadata"`
}

type QueueBindingData struct {
	Msg string `json:"msg"`
}

type QueueBindingMetadata struct {
	Id           string `json:"Id"`
	DequeueCount string `json:"DequeueCount"`
}
