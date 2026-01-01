package usecase

type OrderCreatedEvent struct {
	OrderID     string
	Email       string
	TotalAmount float64
}
