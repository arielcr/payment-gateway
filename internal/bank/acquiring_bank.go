package bank

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/arielcr/payment-gateway/internal/config"
)

type PaymentRequest struct {
	Amount      float64 `json:"amount"`
	CardNumber  string  `json:"card_number"`
	ExpiryMonth string  `json:"expiry_month"`
	ExpiryYear  string  `json:"expiry_year"`
	Cvv         string  `json:"cvv"`
}

type PaymentResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Processor string `json:"processor"`
}

type AcquiringBank struct {
	config config.Application
}

func NewAdquiringBank(config config.Application) AcquiringBank {
	return AcquiringBank{
		config: config,
	}
}

func (a *AcquiringBank) ProcessPayment(request PaymentRequest) (PaymentResponse, error) {
	payloadBytes, err := json.Marshal(request)
	if err != nil {
		return PaymentResponse{}, err
	}

	// Make a POST request to the bank simulator
	resp, err := http.Post(a.config.BankSimulatorHost, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return PaymentResponse{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return PaymentResponse{}, err
	}

	// Parse the response JSON into a PaymentResponse struct
	var response PaymentResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return PaymentResponse{}, err
	}

	return response, nil
}
