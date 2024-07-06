package domain

import "time"

// Session - common project model for all services
type Session struct {
	ID         string
	UserID     string
	Datetime   time.Time
	IP         string
	CountryID  string
	PlatformID int

	UtmContent  string
	UtmTerm     string
	UtmCampaign string
	UtmSource   string
	UtmMedium   string

	TemplateID  int
	TemplateDir string

	IsAbTest      int
	AbTestID      int
	AbTaskID      int
	AbTaskDefault bool
}
