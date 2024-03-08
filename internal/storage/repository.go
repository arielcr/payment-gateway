package storage

import "github.com/arielcr/payment-gateway/internal/models"

type Repository interface {
	Connect() error
	CreatePayment(payment models.Payment) error
}
