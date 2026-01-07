package dto

import (
	"time"

	"github.com/gofrs/uuid"
)

type OrderCreatedEvent struct {
	OrderID     string  `json:"OrderID"`      // PascalCase!
	Email       string  `json:"Email"`        // PascalCase!
	TotalAmount float64 `json:"TotalAmount"`  // PascalCase!
}


type CreateNotificationRequest struct {
	OrderID uuid.UUID `json:"order_id"`
	Type    string    `json:"type"`
	Message string    `json:"message"`
}

type NotificationResponse struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}