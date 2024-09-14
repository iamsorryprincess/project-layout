package domain

import "time"

type Click struct {
	CreatedAt time.Time
	IP        string
}
