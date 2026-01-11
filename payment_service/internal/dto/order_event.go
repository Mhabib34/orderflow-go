package dto

import "github.com/gofrs/uuid"

type OrderCreatedEvent struct {
	OrderID     uuid.UUID `json:"OrderID"`     // PascalCase!
	Email       string    `json:"Email"`       // PascalCase!
	TotalAmount float64   `json:"TotalAmount"` // PascalCase!
}
