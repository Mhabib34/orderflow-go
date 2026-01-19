package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"payment_service/internal/dto"
	"payment_service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentControllerImpl struct {
	PaymentUsecase usecase.PaymentUsecase
}

func NewPaymentController(PaymentUsecase usecase.PaymentUsecase) PaymentController {
	return &PaymentControllerImpl{
		PaymentUsecase: PaymentUsecase,
	}
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
	_, err = controller.PaymentUsecase.CreatePayment(ctx, dto.CreatePaymentRequest{
		OrderID: event.OrderID,
		Amount:  int64(event.TotalAmount),
		Method:  "bank_transfer",
	})

	if err != nil {
		log.Printf("❌ Failed to create payment: %v\n", err)
		return err
	}

	log.Println("✅ Payment created successfully")
	return nil
}

func (c *PaymentControllerImpl) HandleMidtransWebhook(ctx *gin.Context) {
	// ✅ Tambahkan defer recover untuk catch panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("🚨 PANIC in webhook handler: %v\n", r)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}()

	bodyBytes, err := ctx.GetRawData()
	if err != nil {
		log.Println("❌ Failed to read request body:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot read body"})
		return
	}
	
	log.Printf("📦 Webhook Raw Body: %s\n", string(bodyBytes))

	var payload dto.MidtransCallback
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		log.Println("❌ Invalid webhook payload:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	log.Printf("📩 Webhook Parsed: %+v\n", payload)

	// ✅ Pastikan usecase dipanggil
	if err := c.PaymentUsecase.HandleMidtransCallback(ctx, payload); err != nil {
		log.Printf("❌ Callback processing error: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("✅ Webhook processed successfully")
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}