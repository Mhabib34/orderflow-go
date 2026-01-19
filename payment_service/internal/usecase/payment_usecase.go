package usecase

import (
	"context"
	"payment_service/internal/dto"
)

type PaymentUsecase interface {
	CreatePayment(ctx context.Context, request dto.CreatePaymentRequest) (dto.PaymentResponse, error)
	HandleMidtransCallback(ctx context.Context, payload dto.MidtransCallback) error
}


