package models

import "time"

type PaymentResponse struct {
	ID          uint             `json:"id"`
	OrderToken  string           `json:"order_token"`
	Status      PaymentStatus    `json:"status"`
	PaymentInfo PaymentInfo      `json:"payment_info"`
	RedirectUrl string           `json:"redirect_url"`
	Merchant    MerchantResponse `json:"merchant"`
	Customer    CustomerResponse `json:"customer"`
	CreatedAt   time.Time        `json:"created_at"`
}

type PaymentInfo struct {
	Amount      float64     `json:"amount"`
	MethodType  string      `json:"method_type"`
	CardDetails CardDetails `json:"card_details"`
	Processor   string      `json:"processor"`
}

type CardDetails struct {
	CardType       string `json:"card_type"`
	CardBrand      string `json:"card_brand"`
	CardHolder     string `json:"card_holder"`
	LastFourDigits string `json:"last_four_digits"`
}
