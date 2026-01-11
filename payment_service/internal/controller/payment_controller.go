package controller

import "context"

type PaymentController interface {
	Create(ctx context.Context, body []byte) error
}