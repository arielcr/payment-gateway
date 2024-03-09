// Package models provides data models used throughout the application.
package models

import (
	"gorm.io/gorm"
)

// Refund represents a refund entity stored in the database.
type Refund struct {
	gorm.Model         // Embedded gorm.Model for ID, created_at, updated_at, deleted_at fields.
	PaymentID  uint    `gorm:"not null" json:"payment_id" validate:"required"`
	Amount     float64 `gorm:"not null" json:"amount" validate:"required"`
	Reason     string  `gorm:"not null" json:"reason" validate:"required"`
}
