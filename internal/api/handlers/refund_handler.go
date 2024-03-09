package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/arielcr/payment-gateway/internal/bank"
	"github.com/arielcr/payment-gateway/internal/config"
	"github.com/arielcr/payment-gateway/internal/models"
	"github.com/arielcr/payment-gateway/internal/storage"
	"github.com/gin-gonic/gin"
)

type RefundHandler struct {
	store  storage.Repository
	config config.Application
}

func NewRefundHandler(store storage.Repository, config config.Application) *RefundHandler {
	return &RefundHandler{
		store:  store,
		config: config,
	}
}

func (p *RefundHandler) RefundPayment(context *gin.Context) {
	var status models.PaymentStatus
	paymentID := context.Param("paymentID")

	refundRequest := models.RefundRequest{}
	if err := context.BindJSON(&refundRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refundResult, err := p.sendRefundRequest(refundRequest)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if refundResult.Success {
		status = models.Refunded
		_, err := p.createRefund(refundRequest, paymentID, status)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	context.JSON(http.StatusOK, &models.RefundResponse{Status: "refunded"})
}

func (p *RefundHandler) createRefund(refundRequest models.RefundRequest, paymentID string, status models.PaymentStatus) (models.Refund, error) {
	id, err := strconv.Atoi(paymentID)
	if err != nil {
		return models.Refund{}, errors.New("invalid payment id")
	}
	refund := models.Refund{
		Amount:    refundRequest.Amount,
		Reason:    refundRequest.Reason,
		PaymentID: uint(id),
	}

	if err := p.store.CreateRefund(&refund); err != nil {
		return models.Refund{}, err
	}

	if err := p.store.UpdatePaymentStatus(uint(id), status); err != nil {
		return models.Refund{}, err
	}

	return refund, nil
}

func (p *RefundHandler) sendRefundRequest(refundRequest models.RefundRequest) (bank.PaymentResponse, error) {
	request := bank.RefundRequest{
		Amount: refundRequest.Amount,
		Reason: refundRequest.Reason,
	}
	acquiringBank := bank.NewAdquiringBank(p.config)
	response, err := acquiringBank.ProcessRefund(request)
	if err != nil {
		return bank.PaymentResponse{}, err
	}
	return response, nil
}
