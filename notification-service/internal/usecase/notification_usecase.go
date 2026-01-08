package usecase

import (
	"context"
	"notification_service/internal/dto"
	"notification_service/internal/model"
)

type NotificationUsecase interface {
	CreateNotification(ctx context.Context, request dto.CreateNotificationRequest) (dto.NotificationResponse, error)
	GetAll(ctx context.Context, searchRequest dto.SearchNotificationRequest) ([]model.Notifications, int64, error)
}