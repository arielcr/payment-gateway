package handlers

import (
	"net/http"

	"github.com/arielcr/payment-gateway/internal/models"
	"github.com/arielcr/payment-gateway/internal/storage"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	store storage.Repository
}

func NewPaymentHandler(store storage.Repository) *PaymentHandler {
	return &PaymentHandler{
		store: store,
	}
}

func (p *PaymentHandler) ProcessPayment(context *gin.Context) {
	paymentRequest := models.PaymentRequest{}

	if err := context.BindJSON(&paymentRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment := models.Payment{
		OrderToken: paymentRequest.OrderToken,
		MerchantID: paymentRequest.MerchandID,
		Amount:     paymentRequest.Amount,
		Status:     paymentRequest.Status,
		CustomerID: 1,
	}

	if result := p.store.CreatePayment(payment); result != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	context.JSON(http.StatusCreated, &payment)
}
