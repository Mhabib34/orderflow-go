package repository

import (
	"context"
	"notification_service/internal/model"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *model.Notifications) (*model.Notifications, error)
}

