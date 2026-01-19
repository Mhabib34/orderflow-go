package service

type SnapResponse struct {
	Token       string
	RedirectURL string
}

type MidtransService interface {
	CreateSnapPayment(
		paymentID string,
		amount int64,
	) (*SnapResponse, error)
}