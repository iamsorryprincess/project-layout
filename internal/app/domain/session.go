package domain

import "time"

// Session - common project model for all services
type Session struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	Datetime   time.Time `json:"datetime"`
	IP         string    `json:"ip"`
	CountryID  string    `json:"countryId"`
	PlatformID uint8     `json:"platformId"`

	UtmContent  string `json:"utmContent"`
	UtmTerm     string `json:"utmTerm"`
	UtmCampaign string `json:"utmCampaign"`
	UtmSource   string `json:"utmSource"`
	UtmMedium   string `json:"utmMedium"`

	TemplateID  int    `json:"templateId"`
	TemplateDir string `json:"templateDir"`

	IsAbTest      int  `json:"isAbTest"`
	AbTestID      int  `json:"abTestId"`
	AbTaskID      int  `json:"abTaskId"`
	AbTaskDefault bool `json:"abTaskDefault"`
}
