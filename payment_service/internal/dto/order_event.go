package dto

import "github.com/google/uuid"

type OrderCreatedEvent struct {
	OrderID     uuid.UUID `json:"OrderID"`
	Email       string    `json:"Email"`
	TotalAmount float64   `json:"TotalAmount"`
}

type PaymentStatusChangedEvent struct {
	PaymentStatus string `json:"payment_status"`
	PaymentID     uuid.UUID `json:"payment_id"`
	OrderID       uuid.UUID `json:"order_id"`
	PaymentMethod string `json:"payment_method"`
}


type CreatePaymentRequest struct {
	OrderID uuid.UUID `json:"order_id"`
	Amount  int64     `json:"amount"`
	Method  string    `json:"method"`
}

type MidtransCallback struct {
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	FraudStatus       string `json:"fraud_status"`
	PaymentType       string `json:"payment_type"`
}

type PaymentResponse struct {
	ID            uuid.UUID `json:"id"`
	OrderID       uuid.UUID `json:"order_id"`
	Amount        int64     `json:"amount"`
	Method        string    `json:"method"`
	Status        string    `json:"status"`
	PaidAt        string    `json:"paid_at"`
	ExpiredAt     string    `json:"expired_at"`
	ProviderRefId string    `json:"provider_ref_id"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
	PaymentURL    string    `json:"payment_url"`
	PaymentID     string    `json:"payment_id"`
}