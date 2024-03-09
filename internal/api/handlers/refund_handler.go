// Package handlers provides HTTP handlers for processing refund requests.
package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/arielcr/payment-gateway/internal/bank"
	"github.com/arielcr/payment-gateway/internal/config"
	"github.com/arielcr/payment-gateway/internal/models"
	"github.com/arielcr/payment-gateway/internal/storage"
	"github.com/gin-gonic/gin"
)

// RefundHandler handles HTTP requests related to processing refunds.
type RefundHandler struct {
	store  storage.Repository
	config config.Application
	logger *slog.Logger
}

// NewRefundHandler creates a new instance of RefundHandler with the provided store and config.
func NewRefundHandler(store storage.Repository, config config.Application, logger *slog.Logger) *RefundHandler {
	return &RefundHandler{
		store:  store,
		config: config,
		logger: logger,
	}
}

// RefundPayment handles the HTTP POST request to process a refund for a payment.
// It decodes the request body, sends a refund request to the acquiring bank, and updates the payment status.
func (p *RefundHandler) RefundPayment(context *gin.Context) {
	p.logger.Info("Refunding payment")

	var status models.PaymentStatus
	paymentID := context.Param("paymentID")

	refundRequest := models.RefundRequest{}
	if err := context.BindJSON(&refundRequest); err != nil {
		p.logger.Error(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refundResult, err := p.sendRefundRequest(refundRequest)
	if err != nil {
		p.logger.Error(err.Error())
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if refundResult.Success {
		status = models.Refunded
		_, err := p.createRefund(refundRequest, paymentID, status)
		if err != nil {
			p.logger.Error(err.Error())
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	context.JSON(http.StatusOK, &models.RefundResponse{Status: "refunded"})
}

// createRefund creates a refund record in the database and updates the payment status based on the refund request.
func (p *RefundHandler) createRefund(refundRequest models.RefundRequest, paymentID string, status models.PaymentStatus) (models.Refund, error) {
	p.logger.Info("Creating refund")

	id, err := strconv.Atoi(paymentID)
	if err != nil {
		p.logger.Error(err.Error())
		return models.Refund{}, errors.New("invalid payment id")
	}
	refund := models.Refund{
		Amount:    refundRequest.Amount,
		Reason:    refundRequest.Reason,
		PaymentID: uint(id),
	}

	if err := p.store.CreateRefund(&refund); err != nil {
		p.logger.Error(err.Error())
		return models.Refund{}, err
	}

	if err := p.store.UpdatePaymentStatus(uint(id), status); err != nil {
		p.logger.Error(err.Error())
		return models.Refund{}, err
	}

	return refund, nil
}

// sendRefundRequest sends a refund request to the acquiring bank for processing refund.
// It constructs the request using the refund information and application configuration.
func (p *RefundHandler) sendRefundRequest(refundRequest models.RefundRequest) (bank.PaymentResponse, error) {
	p.logger.Info("Sending refund request")

	request := bank.RefundRequest{
		Amount: refundRequest.Amount,
		Reason: refundRequest.Reason,
	}
	acquiringBank := bank.NewAdquiringBank(p.config, p.logger)
	response, err := acquiringBank.ProcessRefund(request)
	if err != nil {
		p.logger.Error(err.Error())
		return bank.PaymentResponse{}, err
	}

	return response, nil
}
