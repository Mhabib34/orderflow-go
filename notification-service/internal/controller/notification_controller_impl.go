package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"notification_service/internal/dto"
	"notification_service/internal/exception"
	"notification_service/internal/helper"
	"notification_service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
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

	// ✅ FIX: Validasi UUID kosong HARUS SEBELUM usecase
	if event.OrderID == uuid.Nil {
		log.Println("❌ OrderID is empty or invalid")
		return fmt.Errorf("order_id is required")
	}


	// Call usecase
	_, err = controller.NotificationUsecase.CreateNotification(ctx, dto.CreateNotificationRequest{
		OrderID: event.OrderID,
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

func (controller *NotificationControllerImpl) GetAll(ctx *gin.Context) {
	page := helper.StringToIntDefault(ctx.Query("page"), 1)
	limit := helper.StringToIntDefault(ctx.Query("limit"), 10)
	isReadParam := ctx.Query("is_read")
	typeParam := ctx.Query("type")

	// ✅ Parse is_read sebagai pointer
	var isRead *bool
	if isReadParam != "" {
		value := helper.StringToBoolDefault(isReadParam, false)
		isRead = &value
	}

	notifications, total, err := controller.NotificationUsecase.GetAll(
		ctx.Request.Context(), 
		dto.SearchNotificationRequest{
			IsRead: isRead,
			Type:   typeParam,
			Limit:  limit,
			Page:   page,
		},
	)
	if err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	webResponse := dto.WebResponse{
		Status: "OK",
		Message: "Orders fetched successfully", 
		Data: notifications,
		Pagination: &dto.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}

	helper.WriteToResponseBody(ctx, http.StatusOK, webResponse)
}

func (controller *NotificationControllerImpl) FindByID(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := helper.StringToUUID(idParam)
	if err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	notification, err := controller.NotificationUsecase.FindByID(ctx.Request.Context(), id)
	if err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	webResponse := dto.WebResponse{
		Status: "OK",
		Message: "Get detail notification successfully", 
		Data: notification,
	}

	helper.WriteToResponseBody(ctx, http.StatusOK, webResponse)
}