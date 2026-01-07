package repository

import (
	"context"
	"notification_service/internal/exception"
	"notification_service/internal/model"

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