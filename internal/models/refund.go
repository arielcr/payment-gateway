package models

import (
	"gorm.io/gorm"
)

type Refund struct {
	gorm.Model
	PaymentID int     `gorm:"not null" json:"payment_id" validate:"required"`
	Amount    float64 `gorm:"not null" json:"amount" validate:"required"`
	Status    int     `gorm:"not null" json:"status" validate:"required"`
}
