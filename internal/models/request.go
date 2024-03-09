package models

type PaymentRequest struct {
	OrderToken    string        `json:"order_token"`
	PaymentSource PaymentSource `json:"payment_source"`
	Amount        float64       `json:"amount"`
	Customer      Customer      `json:"customer"`
	CallbackUrls  CallbackUrls  `json:"callback_urls"`
	MerchandID    uint          `json:"merchand_id"` // se podria mandar como header
}

type PaymentSource struct {
	MethodType string   `json:"method_type"`
	Processor  string   `json:"processor"`
	CardInfo   CardInfo `json:"card_info"`
}

type CardInfo struct {
	CardType        string `json:"card_type"`
	ExpirationMonth string `json:"expiration_month"`
	ExpirationYear  string `json:"expiration_year"`
	CardNumber      string `json:"card_number"`
	CardHolder      string `json:"card_holder"`
	CardCvv         string `json:"card_cvv"`
}

type CallbackUrls struct {
	Success   string `json:"success"`
	Reject    string `json:"reject"`
	Cancelled string `json:"cancelled"`
	Failed    string `json:"failed"`
}
