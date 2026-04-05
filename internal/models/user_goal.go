package models

import (
	"time"
)

type UserGoal struct {
	ID           int       `json:"id,omitempty"`
	UserID       int       `json:"user_id,omitempty"`
	Type         string    `json:"type,omitempty"`
	TargetWeight int       `json:"target_weight,omitempty"`
	WeeklyRate   float64   `json:"weekly_rate,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
