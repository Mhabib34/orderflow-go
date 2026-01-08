package repository

import (
	"context"
	"notification_service/internal/exception"
	"notification_service/internal/model"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type NotificationRepositoryImpl struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &NotificationRepositoryImpl{db: db}
}

func (n *NotificationRepositoryImpl) Create(ctx context.Context, notification *model.Notifications)(*model.Notifications, error) {
	err := n.db.WithContext(ctx).Create(notification).Error
	exception.PanicIfError(err)
	return notification, nil
}

func (n *NotificationRepositoryImpl) GetAll(ctx context.Context, isRead *bool, Type string, limit, page int) ([]model.Notifications, int64, error) {
	var (
		notifications []model.Notifications
		total         int64
	)

	offset := (page - 1) * limit

	// base query
	query := n.db.WithContext(ctx).
		Model(&model.Notifications{})

	// ✅ filter hanya jika type dikirim
	if Type != "" {
		query = query.Where("type = ?", Type)
	}
	
	// ✅ filter hanya jika isRead tidak nil
	if isRead != nil {
		query = query.Where("is_read = ?", *isRead)
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
		Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (n *NotificationRepositoryImpl) FindById(ctx context.Context, id uuid.UUID)(*model.Notifications, error) {
	var notification model.Notifications
	err := n.db.WithContext(ctx).Where("id = ?", id).First(&notification).Error
	exception.PanicIfError(err)
	return &notification, nil
}