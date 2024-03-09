package storage

import "github.com/arielcr/payment-gateway/internal/models"

type Repository interface {
	CreatePayment(payment *models.Payment) error
	CreateCustomer(customer *models.Customer) error
	GetMerchant(merchantID uint) (models.Merchant, error)
	GetCustomer(customerID uint) (models.Customer, error)
	GetPayment(paymentID string) (models.PaymentData, error)
	CreateCreditCard(creditCard *models.CreditCard) error
}
