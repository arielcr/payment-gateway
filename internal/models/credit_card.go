package models

import (
	"gorm.io/gorm"
)

type CreditCard struct {
	gorm.Model
	Token           string `gorm:"not null" json:"token" validate:"required"`
	ExpirationMonth string `gorm:"not null" json:"expiration_date" validate:"required"`
	ExpirationYear  string `gorm:"not null" json:"expiration_year" validate:"required"`
	CardHolder      string `gorm:"not null" json:"card_holder" validate:"required"`
	CardType        string `gorm:"not null" json:"card_type" validate:"required"`
	CardBrand       string `gorm:"not null" json:"card_brand" validate:"required"`
	LastFour        string `gorm:"not null" json:"last_four" validate:"required"`
	CustomerID      uint   `gorm:"not null" json:"customer_id" validate:"required"`
}
