package helper

import (
	"notification_service/internal/dto"
	"notification_service/internal/model"
)

func ToNotificationResponse(notification model.Notifications) dto.NotificationResponse{
	return dto.NotificationResponse{
		ID:        notification.ID,
		OrderID:   notification.OrderID,
		Type:      notification.Type,
		Message:   notification.Message,
		IsRead:    notification.IsRead,
		CreatedAt: notification.CreatedAt,
	}
}