package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"notification_service/internal/dto"
	"notification_service/internal/helper"
	"notification_service/internal/usecase"
)

type NotificationControllerImpl struct {
	NotificationUsecase usecase.NotificationUsecase
}

func NewNotificationController(notificationUsecase usecase.NotificationUsecase) NotificationController {
	return &NotificationControllerImpl{NotificationUsecase: notificationUsecase}
}

func (controller *NotificationControllerImpl) Create(ctx context.Context, body []byte) error {
	// Log raw message untuk debugging
	log.Printf("📦 Raw message: %s\n", string(body))
	
	var event dto.OrderCreatedEvent

	err := json.Unmarshal(body, &event)
	if err != nil {
		log.Printf("❌ Failed to unmarshal JSON: %v\n", err)
		return err
	}

	log.Println("📨 Order Created Event Received")
	log.Printf("   OrderID: %s\n", event.OrderID)
	log.Printf("   Email: %s\n", event.Email)
	log.Printf("   Amount: %.2f\n", event.TotalAmount)

	// Validasi OrderID tidak kosong
	if event.OrderID == "" {
		log.Println("❌ OrderID is empty")
		return fmt.Errorf("order_id is required")
	}

	orderID, err := helper.StringToUUID(event.OrderID)
	if err != nil {
		log.Printf("❌ Invalid UUID format: %v\n", err)
		return err
	}

	// Call usecase
	_, err = controller.NotificationUsecase.CreateNotification(ctx, dto.CreateNotificationRequest{
		OrderID: orderID,
		Type:    "order_created",
		Message: "Your order has been created.",
	})

	if err != nil {
		log.Printf("❌ Failed to create notification: %v\n", err)
		return err
	}

	log.Println("✅ Notification created successfully")
	return nil
}