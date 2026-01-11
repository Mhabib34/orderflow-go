package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"payment_service/internal/dto"

	"github.com/gofrs/uuid"
)

type PaymentControllerImpl struct {
	// PaymentUsecase usecase.PaymentUsecase
}

func NewPaymentController() PaymentController {
	return &PaymentControllerImpl{}
}

func (controller *PaymentControllerImpl) Create(ctx context.Context, body []byte) error {
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

	// ✅ FIX: Validasi UUID kosong HARUS SEBELUM usecase
	if event.OrderID == uuid.Nil {
		log.Println("❌ OrderID is empty or invalid")
		return fmt.Errorf("order_id is required")
	}

	// // Call usecase
	// _, err = controller.NotificationUsecase.CreateNotification(ctx, dto.CreateNotificationRequest{
	// 	OrderID: event.OrderID,
	// 	Type:    "order_created",
	// 	Message: "Your order has been created.",
	// })

	// if err != nil {
	// 	log.Printf("❌ Failed to create notification: %v\n", err)
	// 	return err
	// }

	log.Println("✅ Payment created successfully")
	return nil
}
