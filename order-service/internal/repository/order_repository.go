package repository

import (
	"context"
	"order_service/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Orders)(*model.Orders, error)
}