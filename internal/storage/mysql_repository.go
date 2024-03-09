// Package storage provides interfaces and implementations for interacting with different data storage systems.
package storage

import (
	"errors"
	"fmt"
	"log/slog"
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
	db     *gorm.DB
	logger *slog.Logger
}

// NewMySQLRepository creates a new instance of MySQLRepository.
func NewMySQLRepository(db *gorm.DB, logger *slog.Logger) *MySQLRepository {
	return &MySQLRepository{
		db:     db,
		logger: logger,
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
	m.logger.Info("Creating new payment")

	if result := m.db.Create(&payment); result.Error != nil {
		m.logger.Error(result.Error.Error())
		return result.Error
	}
	return nil
}

// UpdatePaymentStatus updates the status of a payment in the database.
func (m *MySQLRepository) UpdatePaymentStatus(paymentID uint, status models.PaymentStatus) error {
	m.logger.Info("Updating payment status")

	var payment models.Payment
	if err := m.db.First(&payment, paymentID); err.Error != nil {
		m.logger.Error(err.Error.Error())
		return errPaymentNotFound
	}

	payment.Status = status
	if err := m.db.Save(&payment); err.Error != nil {
		m.logger.Error(err.Error.Error())
		return err.Error
	}

	return nil
}

// CreateRefund creates a new refund record in the database.
func (m *MySQLRepository) CreateRefund(refund *models.Refund) error {
	m.logger.Info("Creating new refund")

	var payment models.Payment
	if err := m.db.First(&payment, refund.PaymentID); err.Error != nil {
		m.logger.Error(err.Error.Error())
		return errPaymentNotFound
	}

	// if the refund amount equals zero, then it is a full refund
	if refund.Amount == 0.00 {
		refund.Amount = payment.Amount
	}

	if result := m.db.Create(&refund); result.Error != nil {
		m.logger.Error(result.Error.Error())
		return result.Error
	}
	return nil
}

// CreateCustomer creates a new customer record in the database.
func (m *MySQLRepository) CreateCustomer(customer *models.Customer) error {
	m.logger.Info("Creating new customer")

	if result := m.db.Create(&customer); result.Error != nil {
		m.logger.Error(result.Error.Error())
		return result.Error
	}
	return nil
}

// GetMerchant retrieves a merchant record from the database by ID.
func (m *MySQLRepository) GetMerchant(merchantID uint) (models.Merchant, error) {
	m.logger.Info("Getting a merchant")

	var merchant models.Merchant
	if result := m.db.First(&merchant, merchantID); result.Error != nil {
		m.logger.Error(result.Error.Error())
		return models.Merchant{}, errMerchantNotFound
	}
	return merchant, nil
}

// GetCustomer retrieves a customer record from the database by ID.
func (m *MySQLRepository) GetCustomer(customerID uint) (models.Customer, error) {
	m.logger.Info("Getting a customer")

	var customer models.Customer
	if result := m.db.First(&customer, customerID); result.Error != nil {
		m.logger.Error(result.Error.Error())
		return models.Customer{}, errCustomerNotFound
	}
	return customer, nil
}

// GetPayment retrieves a payment record from the database by ID.
func (m *MySQLRepository) GetPayment(paymentID string) (models.PaymentData, error) {
	m.logger.Info("Getting a payment")

	var payment models.Payment

	id, err := strconv.Atoi(paymentID)
	if err != nil {
		m.logger.Error(errInvalidPaymentId.Error())
		return models.PaymentData{}, errInvalidPaymentId
	}

	if result := m.db.First(&payment, uint(id)); result.Error != nil {
		m.logger.Error(errPaymentNotFound.Error())
		return models.PaymentData{}, errPaymentNotFound
	}

	customer, err := m.GetCustomer(payment.CustomerID)
	if err != nil {
		m.logger.Error(errCustomerNotFound.Error())
		return models.PaymentData{}, errCustomerNotFound
	}

	merchant, err := m.GetMerchant(payment.MerchantID)
	if err != nil {
		m.logger.Error(errMerchantNotFound.Error())
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
	m.logger.Info("Creating new credit card")

	if result := m.db.Create(&creditCard); result.Error != nil {
		m.logger.Error(result.Error.Error())
		return result.Error
	}
	return nil
}
