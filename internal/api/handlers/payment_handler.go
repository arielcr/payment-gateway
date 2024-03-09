// Package handlers provides HTTP handlers for processing payment and refund requests.
package handlers

import (
	"net/http"

	"github.com/arielcr/payment-gateway/internal/bank"
	"github.com/arielcr/payment-gateway/internal/config"
	"github.com/arielcr/payment-gateway/internal/models"
	"github.com/arielcr/payment-gateway/internal/storage"
	"github.com/arielcr/payment-gateway/internal/utils"
	"github.com/gin-gonic/gin"
)

// PaymentHandler handles HTTP requests related to processing payments.
type PaymentHandler struct {
	store  storage.Repository
	config config.Application
}

// NewPaymentHandler creates a new instance of PaymentHandler with the provided store and config.
func NewPaymentHandler(store storage.Repository, config config.Application) *PaymentHandler {
	return &PaymentHandler{
		store:  store,
		config: config,
	}
}

// ProcessPayment handles the HTTP POST request to process a payment.
// It decodes the request body, validates and processes the payment, and sends back a response.
func (p *PaymentHandler) ProcessPayment(context *gin.Context) {
	paymentRequest := models.PaymentRequest{}
	if err := context.BindJSON(&paymentRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	merchant, err := p.getMerchantInfo(paymentRequest.MerchandID)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	customer, err := p.getCustomerInfo(paymentRequest.Customer)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	creditCard, err := p.createCreditCard(paymentRequest.PaymentSource.CardInfo, customer.ID)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactionResult, err := p.sendTransactionRequest(paymentRequest)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := p.createPayment(paymentRequest, transactionResult, customer.ID)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	paymentResponse := p.generateResponse(payment, paymentRequest, merchant, customer, creditCard, transactionResult)

	context.JSON(http.StatusCreated, &paymentResponse)
}

// GetPayment handles the HTTP GET request to retrieve payment information by ID.
// It fetches the payment data from the database and sends back a response.
func (p *PaymentHandler) GetPayment(context *gin.Context) {
	paymentID := context.Param("paymentID")

	paymentData, err := p.store.GetPayment(paymentID)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, &paymentData)
}

// sendTransactionRequest sends a transaction request to the acquiring bank for processing payment.
// It constructs the request using the payment information and application configuration.
func (p *PaymentHandler) sendTransactionRequest(paymentRequest models.PaymentRequest) (bank.PaymentResponse, error) {
	request := bank.PaymentRequest{
		Amount:      paymentRequest.Amount,
		CardNumber:  paymentRequest.PaymentSource.CardInfo.CardNumber,
		ExpiryMonth: paymentRequest.PaymentSource.CardInfo.ExpirationMonth,
		ExpiryYear:  paymentRequest.PaymentSource.CardInfo.ExpirationYear,
		Cvv:         paymentRequest.PaymentSource.CardInfo.CardCvv,
	}
	acquiringBank := bank.NewAdquiringBank(p.config)
	response, err := acquiringBank.ProcessPayment(request)
	if err != nil {
		return bank.PaymentResponse{}, err
	}
	return response, nil
}

// getMerchantInfo retrieves merchant information from the database by ID.
func (p *PaymentHandler) getMerchantInfo(id uint) (models.Merchant, error) {
	merchant, err := p.store.GetMerchant(id)
	if err != nil {
		return models.Merchant{}, err
	}
	return merchant, nil
}

// getCustomerInfo retrieves or creates customer information based on the provided details.
func (p *PaymentHandler) getCustomerInfo(c models.Customer) (models.Customer, error) {
	customer := models.Customer{
		Name:  c.Name,
		Email: c.Email,
	}

	if c.ID == 0 {
		if err := p.store.CreateCustomer(&customer); err != nil {
			return models.Customer{}, err
		}
	} else {
		var err error
		customer, err = p.store.GetCustomer(c.ID)
		if err != nil {
			return models.Customer{}, err
		}
	}

	return customer, nil
}

// createPayment creates a payment record in the database based on the payment request and transaction result.
func (p *PaymentHandler) createPayment(paymentRequest models.PaymentRequest, transactionResult bank.PaymentResponse, customerID uint) (models.Payment, error) {
	var status models.PaymentStatus
	if transactionResult.Success {
		status = models.Succeeded
	} else {
		status = models.Failed
	}
	payment := models.Payment{
		OrderToken: paymentRequest.OrderToken,
		MerchantID: paymentRequest.MerchandID,
		Amount:     paymentRequest.Amount,
		Status:     status,
		CustomerID: customerID,
	}

	if err := p.store.CreatePayment(&payment); err != nil {
		return models.Payment{}, err
	}

	return payment, nil
}

// createCreditCard creates a credit card record in the database based on the provided card information and customer ID.
func (p *PaymentHandler) createCreditCard(cardInfo models.CardInfo, customerID uint) (models.CreditCard, error) {
	err := utils.ValidateCreditCard(cardInfo.CardNumber)
	if err != nil {
		return models.CreditCard{}, err
	}

	creditCardToken, err := utils.TokenizeCreditCard(cardInfo.CardNumber)
	if err != nil {
		return models.CreditCard{}, err
	}

	lastFourDigits, err := utils.GetLastFourDigits(cardInfo.CardNumber)
	if err != nil {
		return models.CreditCard{}, err
	}

	creditCard := models.CreditCard{
		Token:           creditCardToken,
		ExpirationMonth: cardInfo.ExpirationMonth,
		ExpirationYear:  cardInfo.ExpirationYear,
		CardHolder:      cardInfo.CardHolder,
		CardType:        cardInfo.CardType,
		CardBrand:       utils.GetCreditCardBrand(cardInfo.CardNumber),
		LastFour:        lastFourDigits,
		CustomerID:      customerID,
	}

	if err := p.store.CreateCreditCard(&creditCard); err != nil {
		return models.CreditCard{}, err
	}

	return creditCard, nil
}

// generateResponse generates a payment response based on the payment, payment request, merchant, customer, credit card, and transaction result.
func (p *PaymentHandler) generateResponse(
	payment models.Payment,
	paymentRequest models.PaymentRequest,
	merchant models.Merchant,
	customer models.Customer,
	creditCard models.CreditCard,
	transactionResult bank.PaymentResponse) models.PaymentResponse {

	paymentResponse := models.PaymentResponse{
		ID:         payment.ID,
		OrderToken: payment.OrderToken,
		Status:     payment.Status,
		PaymentInfo: models.PaymentInfo{
			Amount:     payment.Amount,
			MethodType: paymentRequest.PaymentSource.MethodType,
			Processor:  transactionResult.Processor,
			CardDetails: models.CardDetails{
				CardType:       paymentRequest.PaymentSource.CardInfo.CardType,
				CardBrand:      creditCard.CardBrand,
				CardHolder:     paymentRequest.PaymentSource.CardInfo.CardHolder,
				LastFourDigits: creditCard.LastFour,
			},
		},
		RedirectUrl: p.getRedirectUrl(paymentRequest, payment.Status),
		Merchant: models.MerchantResponse{
			Name:  merchant.Name,
			Email: merchant.Email,
		},
		Customer: models.CustomerResponse{
			Name:  customer.Name,
			Email: customer.Email,
		},
		CreatedAt: payment.CreatedAt,
	}

	return paymentResponse
}

// getRedirectUrl determines the redirect URL based on the payment status and callback URLs in the payment request.
func (p *PaymentHandler) getRedirectUrl(paymentRequest models.PaymentRequest, status models.PaymentStatus) string {
	var redirectUrl string
	switch status {
	case models.Succeeded:
		redirectUrl = paymentRequest.CallbackUrls.Success
	case models.Failed:
		redirectUrl = paymentRequest.CallbackUrls.Failed
	}
	return redirectUrl
}
