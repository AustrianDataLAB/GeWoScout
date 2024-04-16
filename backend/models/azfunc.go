package models

type InvokeResponse struct {
	Outputs     map[string]interface{}
	Logs        []string
	ReturnValue interface{}
}

type ReturnValue struct {
	Data map[string]interface{}
}
