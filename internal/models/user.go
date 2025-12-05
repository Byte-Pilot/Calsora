package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	Bday      time.Time `json:"bday"`
	CreatedAt time.Time `json:"created_at"`
}
