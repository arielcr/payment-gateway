// Package storage provides interfaces and implementations for interacting with different data storage systems.
package storage

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/arielcr/payment-gateway/internal/config"
	"github.com/arielcr/payment-gateway/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Define custom error messages
var (
	errMerchantNotFound = errors.New("merchant not found")
	errCustomerNotFound = errors.New("customer not found")
	errPaymentNotFound  = errors.New("payment not found")
	errInvalidPaymentId = errors.New("invalid payment id")
)

// MySQLRepository represents a MySQL implementation of the Repository interface.
type MySQLRepository struct {
	db *gorm.DB
}

// NewMySQLRepository creates a new instance of MySQLRepository.
func NewMySQLRepository(db *gorm.DB) *MySQLRepository {
	return &MySQLRepository{
		db: db,
	}
}

// ConnectMySQL connects to MySQL database using the provided configuration parameters.
func ConnectMySQL(config config.RepositoryParameters) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName)

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// CreatePayment creates a new payment record in the database.
func (m *MySQLRepository) CreatePayment(payment *models.Payment) error {
	if result := m.db.Create(&payment); result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdatePaymentStatus updates the status of a payment in the database.
func (m *MySQLRepository) UpdatePaymentStatus(paymentID uint, status models.PaymentStatus) error {
	var payment models.Payment
	if err := m.db.First(&payment, paymentID); err.Error != nil {
		return errPaymentNotFound
	}

	payment.Status = status
	if err := m.db.Save(&payment); err.Error != nil {
		return err.Error
	}

	return nil
}

// CreateRefund creates a new refund record in the database.
func (m *MySQLRepository) CreateRefund(refund *models.Refund) error {
	var payment models.Payment
	if err := m.db.First(&payment, refund.PaymentID); err.Error != nil {
		return errPaymentNotFound
	}

	// if the refund amount equals zero, then it is a full refund
	if refund.Amount == 0.00 {
		refund.Amount = payment.Amount
	}

	if result := m.db.Create(&refund); result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateCustomer creates a new customer record in the database.
func (m *MySQLRepository) CreateCustomer(customer *models.Customer) error {
	if result := m.db.Create(&customer); result.Error != nil {
		return result.Error
	}
	return nil
}

// GetMerchant retrieves a merchant record from the database by ID.
func (m *MySQLRepository) GetMerchant(merchantID uint) (models.Merchant, error) {
	var merchant models.Merchant
	if result := m.db.First(&merchant, merchantID); result.Error != nil {
		return models.Merchant{}, errMerchantNotFound
	}
	return merchant, nil
}

// GetCustomer retrieves a customer record from the database by ID.
func (m *MySQLRepository) GetCustomer(customerID uint) (models.Customer, error) {
	var customer models.Customer
	if result := m.db.First(&customer, customerID); result.Error != nil {
		return models.Customer{}, errCustomerNotFound
	}
	return customer, nil
}

// GetPayment retrieves a payment record from the database by ID.
func (m *MySQLRepository) GetPayment(paymentID string) (models.PaymentData, error) {
	var payment models.Payment

	id, err := strconv.Atoi(paymentID)
	if err != nil {
		return models.PaymentData{}, errInvalidPaymentId
	}

	if result := m.db.First(&payment, uint(id)); result.Error != nil {
		return models.PaymentData{}, errPaymentNotFound
	}

	customer, err := m.GetCustomer(payment.CustomerID)
	if err != nil {
		return models.PaymentData{}, errCustomerNotFound
	}

	merchant, err := m.GetMerchant(payment.MerchantID)
	if err != nil {
		return models.PaymentData{}, errMerchantNotFound
	}

	paymentData := models.PaymentData{
		OrderToken: payment.OrderToken,
		Amount:     payment.Amount,
		Status:     payment.Status,
		Customer: models.CustomerResponse{
			Name:  customer.Name,
			Email: customer.Email,
		},
		Merchant: models.MerchantResponse{
			Name:  merchant.Name,
			Email: merchant.Email,
		},
		CreatedAt: payment.CreatedAt,
	}

	return paymentData, nil
}

// CreateCreditCard creates a new credit card record in the database.
func (m *MySQLRepository) CreateCreditCard(creditCard *models.CreditCard) error {
	if result := m.db.Create(&creditCard); result.Error != nil {
		return result.Error
	}
	return nil
}
