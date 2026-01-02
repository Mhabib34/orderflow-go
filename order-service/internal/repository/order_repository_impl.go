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

func (o *OrderRepositoryImpl) GetAll(
	ctx context.Context,
	status string,
	limit, page int,
) ([]model.Orders, int64, error) {

	var (
		orders []model.Orders
		total  int64
	)

	offset := (page - 1) * limit

	// base query
	query := o.db.WithContext(ctx).
		Model(&model.Orders{})

	// ✅ filter hanya jika status dikirim
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ambil data
	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}
