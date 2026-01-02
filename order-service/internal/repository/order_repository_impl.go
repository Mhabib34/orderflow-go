package repository

import (
	"context"
	"order_service/internal/exception"
	"order_service/internal/model"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type OrderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &OrderRepositoryImpl{db: db}
}

func (o *OrderRepositoryImpl) CreateOrder(ctx context.Context, order *model.Orders)(*model.Orders, error){
	err := o.db.WithContext(ctx).Create(order).Error
	exception.PanicIfError(err)
	return order, nil
}

func (o *OrderRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID)(*model.Orders, error) {
	var order model.Orders
	err := o.db.WithContext(ctx).Where("id = ?", id).First(&order).Error
	exception.PanicIfError(err)
	return &order, nil
}