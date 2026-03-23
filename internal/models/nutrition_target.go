package models

import "time"

type NutritionTarget struct {
	ID        int       `json:"id,omitempty"`
	UserID    int       `json:"user_id,omitempty"`
	Cal       int       `json:"calories,omitempty"`
	Protein   float64   `json:"protein,omitempty"`
	Carbs     float64   `json:"carbs,omitempty"`
	Fats      float64   `json:"fat,omitempty"`
	IsCustom  bool      `json:"is_custom"`
	CreatedAt time.Time `json:"created_at"`
}
