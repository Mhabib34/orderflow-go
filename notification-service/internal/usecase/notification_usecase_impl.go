package usecase

import (
	"context"
	"notification_service/internal/dto"
	"notification_service/internal/exception"
	"notification_service/internal/helper"
	"notification_service/internal/model"
	"notification_service/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
)

type NotificationUsecaseImpl struct {
	NotificationRepository repository.NotificationRepository
	Validate       *validator.Validate
}

func NewNotificationUsecase(notificationRepository repository.NotificationRepository, validate *validator.Validate) NotificationUsecase {
	return &NotificationUsecaseImpl{NotificationRepository: notificationRepository, Validate: validate}
}

func (service *NotificationUsecaseImpl) CreateNotification(ctx context.Context, request dto.CreateNotificationRequest) (dto.NotificationResponse, error) {
	err := service.Validate.Struct(request)
	exception.PanicIfError(err)

	notification := &model.Notifications{
		OrderID: request.OrderID,
		Type:    request.Type,
		Message: request.Message,
	}

	notification, err = service.NotificationRepository.Create(ctx, notification)
	exception.PanicIfError(err)

	return helper.ToNotificationResponse(*notification), nil
}

func (service *NotificationUsecaseImpl) GetAll(ctx context.Context, searchRequest dto.SearchNotificationRequest) ([]model.Notifications, int64, error) {
	notifications, total, err := service.NotificationRepository.GetAll(
		ctx, 
		searchRequest.IsRead, 
		searchRequest.Type, 
		searchRequest.Limit, 
		searchRequest.Page,
	)
	exception.PanicIfError(err)

	return notifications, total, nil
}

func (service *NotificationUsecaseImpl) FindByID(ctx context.Context, id uuid.UUID) (dto.NotificationResponse, error) {
	notification, err := service.NotificationRepository.FindById(ctx, id)
	exception.PanicIfError(err)

	return helper.ToNotificationResponse(*notification), nil
}