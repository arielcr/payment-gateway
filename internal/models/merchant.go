package models

import (
	"gorm.io/gorm"
)

type Merchant struct {
	gorm.Model
	ApiToken string `gorm:"not null" json:"api_token" validate:"required"`
	Name     string `gorm:"not null" json:"name" validate:"required"`
	Email    string `gorm:"not null" json:"email" validate:"required"`
	Country  string `json:"country"`
	Address  string `json:"address"`
}
