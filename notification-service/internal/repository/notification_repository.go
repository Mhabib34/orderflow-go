package repository

import (
	"context"
	"notification_service/internal/model"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *model.Notifications) (*model.Notifications, error)
	GetAll(ctx context.Context, isRead *bool, Type string , limit, page int) ([]model.Notifications, int64, error)
}

