package usecase

import (
	"context"
	"order_service/internal/dto"

	"github.com/gofrs/uuid"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, request dto.CreateOrderRequest) (dto.OrderResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (dto.OrderResponse, error)
}