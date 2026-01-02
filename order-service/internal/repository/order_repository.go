package repository

import (
	"context"
	"order_service/internal/model"

	"github.com/gofrs/uuid"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Orders)(*model.Orders, error)
	FindByID(ctx context.Context, id uuid.UUID)(*model.Orders, error)
}