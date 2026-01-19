package repository

import (
	"context"
	"fmt"
	"log"
	"payment_service/internal/exception"
	"payment_service/internal/model"

	"gorm.io/gorm"
)

type PaymentRepositoryImpl struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &PaymentRepositoryImpl{
		db: db,
	}
}

func (repository *PaymentRepositoryImpl) Create(ctx context.Context, payment *model.Payments) (*model.Payments, error) {
	err := repository.db.WithContext(ctx).Create(payment).Error
	exception.PanicIfError(err)
	return payment, nil
}

func (r *PaymentRepositoryImpl) UpdateStatusByPaymentID(
	ctx context.Context,
	paymentID string,
	status string,
) error {
	result := r.db.
		WithContext(ctx).
		Model(&model.Payments{}).
		Where("payment_id = ?", paymentID).
		Update("status", status)

	if result.Error != nil {
		log.Printf("❌ DB Error: %v\n", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Printf("⚠️ No payment found with payment_id: %s\n", paymentID)
		return fmt.Errorf("payment with payment_id %s not found", paymentID)
	}

	log.Printf("✅ Updated %d row(s) for payment_id: %s\n", result.RowsAffected, paymentID)
	return nil
}

func (r *PaymentRepositoryImpl) FindByPaymentID(ctx context.Context, paymentID string) (*model.Payments, error) {
	var payment model.Payments
	
	err := r.db.WithContext(ctx).
		Where("payment_id = ?", paymentID).
		First(&payment).Error
	
	if err != nil {
		return nil, err
	}
	
	return &payment, nil
}