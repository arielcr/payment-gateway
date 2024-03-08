package models

type PaymentRequest struct {
	OrderToken string   `json:"order_token"`
	Amount     float64  `json:"amount"`
	Status     int      `json:"status"`
	Customer   Customer `json:"customer"`
	MerchandID int      `json:"merchand_id"`
}
