package repository

import (
	"context"
	"payment_service/internal/model"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *model.Payments) (*model.Payments, error)
	UpdateStatusByPaymentID(ctx context.Context, paymentID string, status string) error
}