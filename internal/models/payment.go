// Package models provides data models used throughout the application.
package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// PaymentStatus represents the status of a payment.
type PaymentStatus int

const (
	Pending PaymentStatus = iota + 1
	Succeeded
	Failed
	Cancelled
	Refunded
	Processed
	Authorized
)

// Payment represents a payment entity stored in the database.
type Payment struct {
	gorm.Model               // Embedded gorm.Model for ID, created_at, updated_at, deleted_at fields.
	OrderToken string        `gorm:"not null" json:"order_token" validate:"required"`
	CustomerID uint          `gorm:"not null" json:"customer_id"`
	MerchantID uint          `gorm:"not null" json:"merchant_id" validate:"required"`
	Amount     float64       `gorm:"not null" json:"amount" validate:"required"`
	Status     PaymentStatus `gorm:"not null" json:"status" validate:"required"`
}

// PaymentData represents data for a payment used in responses.
type PaymentData struct {
	OrderToken string           `json:"order_token"`
	Amount     float64          `json:"amount"`
	Status     PaymentStatus    `json:"status"`
	CreatedAt  time.Time        `json:"created_at"`
	Customer   CustomerResponse `json:"customer"`
	Merchant   MerchantResponse `json:"merchant"`
}

// String converts a PaymentStatus to its string representation.
func (s PaymentStatus) String() string {
	switch s {
	case Pending:
		return "pending"
	case Succeeded:
		return "succeeded"
	case Failed:
		return "failed"
	case Cancelled:
		return "cancelled"
	case Refunded:
		return "refunded"
	case Processed:
		return "processed"
	case Authorized:
		return "authorized"
	default:
		return "unknown"
	}
}

// MarshalJSON marshals a PaymentStatus to JSON.
func (s PaymentStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s.String())), nil
}

// Scan scans a value into a PaymentStatus.
func (ps *PaymentStatus) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*ps = ps.ConvertStringToPaymentStatus(string(v))
		return nil
	default:
		return fmt.Errorf("unsupported Scan type: %T", value)
	}
}

// ConvertStringToPaymentStatus converts a string to a PaymentStatus.
func (ps *PaymentStatus) ConvertStringToPaymentStatus(status string) PaymentStatus {
	switch status {
	case "pending":
		return Pending
	case "succeeded":
		return Succeeded
	case "failed":
		return Failed
	case "cancelled":
		return Cancelled
	case "refunded":
		return Refunded
	case "processed":
		return Processed
	case "authorized":
		return Authorized
	default:
		return 0
	}
}
