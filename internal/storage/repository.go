package storage

import "github.com/arielcr/payment-gateway/internal/models"

type Repository interface {
	CreatePayment(payment *models.Payment) error
	CreateRefund(refund *models.Refund) error
	CreateCustomer(customer *models.Customer) error
	CreateCreditCard(creditCard *models.CreditCard) error
	GetMerchant(merchantID uint) (models.Merchant, error)
	GetCustomer(customerID uint) (models.Customer, error)
	GetPayment(paymentID string) (models.PaymentData, error)
	UpdatePaymentStatus(paymentID uint, status models.PaymentStatus) error
}
