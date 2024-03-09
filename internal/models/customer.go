// Package models provides data models used throughout the application.
package models

import (
	"gorm.io/gorm"
)

// Customer represents a customer entity stored in the database.
type Customer struct {
	gorm.Model        // Embedded gorm.Model for ID, created_at, updated_at, deleted_at fields.
	Name       string `gorm:"not null" json:"name" validate:"required"`
	Email      string `gorm:"not null" json:"email" validate:"required"`
}

// CustomerResponse represents a response model for customer data.
type CustomerResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
