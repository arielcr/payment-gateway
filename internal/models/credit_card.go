package models

import (
	"gorm.io/gorm"
)

type CreditCard struct {
	gorm.Model
	Token          string `gorm:"not null" json:"token" validate:"required"`
	ExpirationDate string `gorm:"not null" json:"expiration_date" validate:"required"`
	CardType       string `gorm:"not null" json:"card_type" validate:"required"`
	CustomerID     int    `gorm:"not null" json:"customer_id" validate:"required"`
	Status         int    `gorm:"not null" json:"status" validate:"required"`
}
