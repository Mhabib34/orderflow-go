package usecase

import (
	"context"
	"notification_service/internal/dto"
)

type NotificationUsecase interface {
	CreateNotification(ctx context.Context, request dto.CreateNotificationRequest) (dto.NotificationResponse, error)
}