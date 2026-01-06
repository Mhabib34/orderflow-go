package dto

type OrderCreatedEvent struct {
	OrderID     string  `json:"order_id"`
	Email       string  `json:"email"`
	TotalAmount float64 `json:"total_amount"`
}
