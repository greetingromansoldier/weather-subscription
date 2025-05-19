package models

type Weather struct {
	temperature float64 `json:"temperature"`
	humidity    float64 `json:"humidity"`
	description string  `json:"description"`
}
