// Package models provides data models used throughout the application.
package models

import (
	"gorm.io/gorm"
)

// CreditCard represents a credit card entity stored in the database.
type CreditCard struct {
	gorm.Model             // Embedded gorm.Model for ID, created_at, updated_at, deleted_at fields.
	Token           string `gorm:"not null" json:"token" validate:"required"`
	ExpirationMonth string `gorm:"not null" json:"expiration_date" validate:"required"`
	ExpirationYear  string `gorm:"not null" json:"expiration_year" validate:"required"`
	CardHolder      string `gorm:"not null" json:"card_holder" validate:"required"`
	CardType        string `gorm:"not null" json:"card_type" validate:"required"`
	CardBrand       string `gorm:"not null" json:"card_brand" validate:"required"`
	LastFour        string `gorm:"not null" json:"last_four" validate:"required"`
	CustomerID      uint   `gorm:"not null" json:"customer_id" validate:"required"`
}
