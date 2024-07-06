package domain

// Event - common project model for all services
type Event struct {
	IP         string `json:"ip"`
	CountryID  string `json:"countryId"`
	PlatformID int    `json:"platformId"`
}
