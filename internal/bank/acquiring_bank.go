// Package bank provides functionality for interacting with the acquiring bank's API.
package bank

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/arielcr/payment-gateway/internal/config"
)

// PaymentRequest represents a payment request sent to the acquiring bank.
type PaymentRequest struct {
	Amount      float64 `json:"amount"`
	CardNumber  string  `json:"card_number"`
	ExpiryMonth string  `json:"expiry_month"`
	ExpiryYear  string  `json:"expiry_year"`
	Cvv         string  `json:"cvv"`
}

// RefundRequest represents a refund request sent to the acquiring bank.
type RefundRequest struct {
	Amount float64 `json:"amount"`
	Reason string  `json:"card_number"`
}

// PaymentResponse represents the response received from the acquiring bank for a payment or refund request.
type PaymentResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Processor string `json:"processor"`
}

// AcquiringBank manages interactions with the acquiring bank's API.
type AcquiringBank struct {
	config config.Application
}

// NewAdquiringBank creates a new instance of AcquiringBank with the provided configuration.
func NewAdquiringBank(config config.Application) AcquiringBank {
	return AcquiringBank{
		config: config,
	}
}

// ProcessPayment sends a payment request to the acquiring bank's API and returns the response.
func (a *AcquiringBank) ProcessPayment(request PaymentRequest) (PaymentResponse, error) {
	payloadBytes, err := json.Marshal(request)
	if err != nil {
		return PaymentResponse{}, err
	}

	resp, err := http.Post(a.config.BankSimulatorHost+"/process", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return PaymentResponse{}, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return PaymentResponse{}, err
	}

	var response PaymentResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return PaymentResponse{}, err
	}

	return response, nil
}

// ProcessRefund sends a refund request to the acquiring bank's API and returns the response.
func (a *AcquiringBank) ProcessRefund(request RefundRequest) (PaymentResponse, error) {
	payloadBytes, err := json.Marshal(request)
	if err != nil {
		return PaymentResponse{}, err
	}

	resp, err := http.Post(a.config.BankSimulatorHost+"/refund", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return PaymentResponse{}, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return PaymentResponse{}, err
	}

	var response PaymentResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return PaymentResponse{}, err
	}

	return response, nil
}
