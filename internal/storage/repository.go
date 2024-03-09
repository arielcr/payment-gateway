// Package storage provides interfaces and implementations for interacting with different data storage systems.
package storage

import "github.com/arielcr/payment-gateway/internal/models"

// Repository defines the interface for interacting with the storage system.
type Repository interface {
	// CreatePayment creates a new payment record in the storage system.
	CreatePayment(payment *models.Payment) error

	// CreateRefund creates a new refund record in the storage system.
	CreateRefund(refund *models.Refund) error

	// CreateCustomer creates a new customer record in the storage system.
	CreateCustomer(customer *models.Customer) error

	// CreateCreditCard creates a new credit card record in the storage system.
	CreateCreditCard(creditCard *models.CreditCard) error

	// GetMerchant retrieves a merchant record from the storage system by ID.
	GetMerchant(merchantID uint) (models.Merchant, error)

	// GetCustomer retrieves a customer record from the storage system by ID.
	GetCustomer(customerID uint) (models.Customer, error)

	// GetPayment retrieves a payment record from the storage system by ID.
	GetPayment(paymentID string) (models.PaymentData, error)

	// UpdatePaymentStatus updates the status of a payment in the storage system.
	UpdatePaymentStatus(paymentID uint, status models.PaymentStatus) error
}
