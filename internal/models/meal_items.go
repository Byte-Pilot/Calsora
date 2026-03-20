package models

import "time"

type MealItems struct {
	ID         int       `json:"id"`
	MealID     int       `json:"meal_id"`
	Name       string    `json:"name"`
	Grams      int       `json:"grams"`
	Cal        int       `json:"cal"`
	Protein    float64   `json:"protein"`
	Carbs      float64   `json:"carbs"`
	Fats       float64   `json:"fats"`
	Confidence float64   `json:"confidence"`
	CreatedAt  time.Time `json:"created_at"`
}
