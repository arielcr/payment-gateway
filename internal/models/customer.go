package models

import (
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	Name  string `gorm:"not null" json:"name" validate:"required"`
	Email string `gorm:"not null" json:"email" validate:"required"`
}
