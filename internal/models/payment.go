package models

import (
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	OrderToken string  `gorm:"not null" json:"order_token" validate:"required"`
	CustomerID int     `gorm:"not null" json:"customer_id" validate:"required"`
	MerchantID int     `gorm:"not null" json:"merchant_id" validate:"required"`
	Amount     float64 `gorm:"not null" json:"amount" validate:"required"`
	Status     int     `gorm:"not null" json:"status" validate:"required"`
}
