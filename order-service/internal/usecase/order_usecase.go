package usecase

import (
	"context"
	"order_service/internal/dto"
	"order_service/internal/model"

	"github.com/gofrs/uuid"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, request dto.CreateOrderRequest) (dto.OrderResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (dto.OrderResponse, error)
	GetAll(ctx context.Context, status string, limit, page int) ([]model.Orders, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, request dto.UpdateStatusRequest) (dto.OrderResponse, error)
	UpdateOrderStatus(ctx context.Context, event dto.PaymentStatusChangedEvent) error
}