package models

import (
	"time"
)

type UserProfile struct {
	UserID    int       `json:"user_id,omitempty"`
	Sex       string    `json:"sex,omitempty"`
	Height    int       `json:"height,omitempty"`
	Weight    int       `json:"weight,omitempty"`
	Activity  int       `json:"activity,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}
