package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

type PaymentRequest struct {
	Amount      float64 `json:"amount"`
	CardNumber  string  `json:"card_number"`
	ExpiryMonth string  `json:"expiry_month"`
	ExpiryYear  string  `json:"expiry_year"`
	CVV         string  `json:"cvv"`
}

type RefundRequest struct {
	Amount float64 `json:"amount"`
	Reason string  `json:"reason"`
}

type PaymentResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Processor string `json:"processor"`
}

func handlePaymentRequest(w http.ResponseWriter, r *http.Request) {
	var paymentRequest PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentRequest); err != nil {
		fmt.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	success := rand.Intn(2) == 0
	var message string
	if success {
		message = "Payment succeeded"
	} else {
		message = "Payment failed"
	}

	paymentResponse := PaymentResponse{
		Success:   success,
		Message:   message,
		Processor: "Awesome Bank",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(paymentResponse); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func handleRefundRequest(w http.ResponseWriter, r *http.Request) {
	var refundRequest RefundRequest
	if err := json.NewDecoder(r.Body).Decode(&refundRequest); err != nil {
		fmt.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	success := true
	message := "Refund succeeded"

	paymentResponse := PaymentResponse{
		Success:   success,
		Message:   message,
		Processor: "Awesome Bank",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(paymentResponse); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/payment/process", handlePaymentRequest)
	http.HandleFunc("/payment/refund", handleRefundRequest)

	log.Println("Bank Simulator started on port 8090")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
