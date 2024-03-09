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

type PaymentHandler struct {
	store  storage.Repository
	config config.Application
}

func NewPaymentHandler(store storage.Repository, config config.Application) *PaymentHandler {
	return &PaymentHandler{
		store:  store,
		config: config,
	}
}

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

func (p *PaymentHandler) GetPayment(context *gin.Context) {
	paymentID := context.Param("paymentID")

	paymentData, err := p.store.GetPayment(paymentID)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, &paymentData)
}

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

func (p *PaymentHandler) getMerchantInfo(id uint) (models.Merchant, error) {
	merchant, err := p.store.GetMerchant(id)
	if err != nil {
		return models.Merchant{}, err
	}
	return merchant, nil
}

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
