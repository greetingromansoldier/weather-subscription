package models

type Weather struct {
	// City should not be here btw
	City        string  `json:"city,omitempty"` //show in json if provided
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Description string  `json:"description"`
}
