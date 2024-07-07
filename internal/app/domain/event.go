package domain

import "time"

// Event - common project model for all services
type Event struct {
	CreatedAt  time.Time `json:"createdAt"`
	IP         string    `json:"ip"`
	CountryID  string    `json:"countryId"`
	PlatformID int       `json:"platformId"`
}
