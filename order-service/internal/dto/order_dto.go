package dto

import (
	"time"

	"github.com/gofrs/uuid"
)

type CreateOrderRequest struct {
	TotalAmount float64 `json:"total_amount" validate:"required"`
	Email       string  `json:"email" validate:"required"`
}

type OrderResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
}