package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Id         uint      `gorm:"PRIMARYKEY"`
	Email      string    `gorm:"type:varchar(255);unique"`
	Password   string    `gorm:"type:varchar(255)"`
	Bday       time.Time `gorm:"type:date"`
	Created_at time.Time `gorm:"type:timestamp"`
}
