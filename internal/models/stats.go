package models

import "time"

type Stats struct {
	Day     time.Time `json:"day"`
	Cal     int       `json:"cal"`
	Protein float64   `json:"protein"`
	Carbs   float64   `json:"carbs"`
	Fats    float64   `json:"fats"`
}
