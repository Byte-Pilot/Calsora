package models

import "time"

type Subscriptions struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Plan      string    `json:"plan"`
	Status    string    `json:"status"`
	StartedAT time.Time `json:"started_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
