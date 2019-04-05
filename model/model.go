package model

// Rates is the model to unmarshal the response body
type Rates struct {
	Base  string
	Rates map[string]float32
	Date  string
}
