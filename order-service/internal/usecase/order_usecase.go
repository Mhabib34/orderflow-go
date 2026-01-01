package usecase

import (
	"context"
	"order_service/internal/dto"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, request dto.CreateOrderRequest) (dto.OrderResponse, error)
}