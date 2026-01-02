package repository

import (
	"context"
	"order_service/internal/model"

	"github.com/gofrs/uuid"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Orders)(*model.Orders, error)
	FindByID(ctx context.Context, id uuid.UUID)(*model.Orders, error)
	GetAll(ctx context.Context, status string, limit, page int) ([]model.Orders, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, order *model.Orders)(*model.Orders, error)
}