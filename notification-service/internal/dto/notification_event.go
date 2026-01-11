package dto

import (
	"time"

	"github.com/gofrs/uuid"
)


type CreateNotificationRequest struct {
	OrderID uuid.UUID `json:"order_id"`
	Type    string    `json:"type"`
	Message string    `json:"message"`
}

type SearchNotificationRequest struct {
	IsRead *bool  `json:"is_read"`
	Type   string `json:"type"`
	Limit  int    `json:"limit"`
	Page   int    `json:"page"`
}

type MarkNotificationAsReadRequest struct {
	IsRead bool `validate:"required" json:"is_read"`
}

type NotificationResponse struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}