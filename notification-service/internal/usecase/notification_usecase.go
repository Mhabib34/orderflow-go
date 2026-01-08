package usecase

import (
	"context"
	"notification_service/internal/dto"
	"notification_service/internal/model"

	"github.com/gofrs/uuid"
)

type NotificationUsecase interface {
	CreateNotification(ctx context.Context, request dto.CreateNotificationRequest) (dto.NotificationResponse, error)
	GetAll(ctx context.Context, searchRequest dto.SearchNotificationRequest) ([]model.Notifications, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (dto.NotificationResponse, error)
	Update(ctx context.Context, id uuid.UUID, request dto.MarkNotificationAsReadRequest) error
}