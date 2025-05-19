package models

// handle enum for frequency
type FrequencyType string

const (
	Hourly FrequencyType = "hourly"
	Daily  FrequencyType = "daily"
)

type Subscription struct {
	email     string        `json:"email"`
	city      string        `json:"city"`
	frequency FrequencyType `json:"frequency"`
	confirmed bool          `json:"confirmed"`
}
