package models

// handle enum for frequency
type FrequencyType string

const (
	Hourly FrequencyType = "hourly"
	Daily  FrequencyType = "daily"
)

type Subscription struct {
	Email     string        `json:"email"`
	City      string        `json:"city"`
	Frequency FrequencyType `json:"frequency" validate:"required,oneof=hourly daily"`
	Confirmed bool          `json:"confirmed"`
}
