// Package models provides data models used throughout the application.
package models

import (
	"gorm.io/gorm"
)

// Merchant represents a merchant entity stored in the database.
type Merchant struct {
	gorm.Model         // Embedded gorm.Model for ID, created_at, updated_at, deleted_at fields.
	ApiToken    string `gorm:"not null" json:"api_token" validate:"required"`
	Name        string `gorm:"not null" json:"name" validate:"required"`
	Email       string `gorm:"not null" json:"email" validate:"required"`
	Country     string `json:"country"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

// MerchantResponse represents a response model for merchant data.
type MerchantResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
