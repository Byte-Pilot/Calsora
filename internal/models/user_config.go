package models

import (
	"time"
)

type UserConfig struct {
	UserID    uint      `json:"user_id,omitempty"`
	Sex       string    `json:"sex,omitempty"`
	Height    uint      `json:"height,omitempty"`
	Weight    uint      `json:"weight,omitempty"`
	CalGoal   uint      `json:"cal_goal,omitempty"`
	Activity  uint      `json:"activity,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}
