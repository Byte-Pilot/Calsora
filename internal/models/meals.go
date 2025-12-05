package models

import "time"

type Meals struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Cal       int       `json:"cal"`
	Protein   float32   `json:"protein"`
	Carbs     float32   `json:"carbs"`
	Fats      float32   `json:"fats"`
	CreatedAt time.Time `json:"created_at"`
}
