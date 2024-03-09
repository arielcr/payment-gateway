package models

import (
	"gorm.io/gorm"
)

type Refund struct {
	gorm.Model
	PaymentID uint    `gorm:"not null" json:"payment_id" validate:"required"`
	Amount    float64 `gorm:"not null" json:"amount" validate:"required"`
	Reason    string  `gorm:"not null" json:"reason" validate:"required"`
}
