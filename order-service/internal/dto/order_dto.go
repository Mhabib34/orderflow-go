package dto

import (
	"time"

	"github.com/gofrs/uuid"
)

type OrderStatus string

const (
	OrderStatusPending  OrderStatus = "pending"
	OrderStatusPaid     OrderStatus = "paid"
	OrderStatusCanceled OrderStatus = "canceled"
)

type CreateOrderRequest struct {
	TotalAmount float64 `json:"total_amount" validate:"required"`
	Email       string  `json:"email" validate:"required"`
}

type UpdateStatusRequest struct {
	Status OrderStatus `json:"status" validate:"required,oneof=pending paid canceled"`
}

type OrderResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
}