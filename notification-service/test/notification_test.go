package test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"notification_service/internal/controller"
	"notification_service/internal/dto"
	"notification_service/internal/model"
	"notification_service/internal/repository"
	"notification_service/internal/usecase"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	testDB         *gorm.DB
	testController controller.NotificationController
)

func setupTestDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=habib123 dbname=notification_service_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&model.Notifications{})
	if err != nil {
		panic(err)
	}

	return db
}

func setupController(db *gorm.DB) controller.NotificationController {
	validate := validator.New()
	repo := repository.NewNotificationRepository(db)
	usecase := usecase.NewNotificationUsecase(repo, validate)
	controller := controller.NewNotificationController(usecase)

	return controller
}

func truncateNotifications(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE notifications")
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	testDB = setupTestDB()
	testController = setupController(testDB)

	code := m.Run()

	os.Exit(code)
}

// ========== TEST RABBITMQ MESSAGE PROCESSING ==========

func TestProcessOrderCreatedMessageSuccess(t *testing.T) {
	truncateNotifications(testDB)

	// ===== simulate RabbitMQ message =====
	orderID := uuid.Must(uuid.NewV4())
	message := dto.OrderCreatedEvent{
		OrderID:     orderID,
		Email:       "test@example.com",
		TotalAmount: 250000,
	}

	messageBytes, err := json.Marshal(message)
	assert.Nil(t, err)

	// ===== process message =====
	ctx := context.Background()
	err = testController.Create(ctx, messageBytes)

	// ===== assert no error =====
	assert.Nil(t, err)

	// ===== assert database =====
	var notification model.Notifications
	err = testDB.First(&notification, "order_id = ?", orderID).Error
	assert.Nil(t, err)

	assert.Equal(t, orderID.String(), notification.OrderID.String())
	assert.Equal(t, "order_created", notification.Type)
	assert.Equal(t, "Your order has been created.", notification.Message)
	assert.False(t, notification.IsRead)
}

func TestProcessOrderCreatedMessageWithInvalidJSON(t *testing.T) {
	truncateNotifications(testDB)

	// ===== invalid JSON =====
	invalidJSON := []byte(`{"OrderID": "invalid-uuid", "Email": }`)

	// ===== process message =====
	ctx := context.Background()
	err := testController.Create(ctx, invalidJSON)

	// ===== assert error =====
	assert.NotNil(t, err)

	// ===== assert no record in database =====
	var count int64
	testDB.Model(&model.Notifications{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestProcessOrderCreatedMessageWithEmptyOrderID(t *testing.T) {
	truncateNotifications(testDB)

	// ===== message with empty OrderID =====
	message := dto.OrderCreatedEvent{
		OrderID:     uuid.Nil,
		Email:       "test@example.com",
		TotalAmount: 250000,
	}

	messageBytes, err := json.Marshal(message)
	assert.Nil(t, err)

	// ===== process message =====
	ctx := context.Background()
	err = testController.Create(ctx, messageBytes)

	// ===== assert error =====
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "order_id is required")

	// ===== assert no record in database =====
	var count int64
	testDB.Model(&model.Notifications{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestProcessOrderCreatedMessageWithInvalidUUID(t *testing.T) {
	truncateNotifications(testDB)

	// ===== message with invalid UUID =====
	message := dto.OrderCreatedEvent{
		OrderID:     uuid.Must(uuid.NewV4()),
		Email:       "test@example.com",
		TotalAmount: 250000,
	}

	_, err := json.Marshal(message)
	assert.Nil(t, err)

	// ===== process message =====
	ctx := context.Background()
	err = testController.Create(ctx, []byte(`{"OrderID":"not-a-valid-uuid","Email":"test@example.com","TotalAmount":250000}`))

	// ===== assert error =====
	assert.NotNil(t, err)

	// ===== assert no record in database =====
	var count int64
	testDB.Model(&model.Notifications{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestProcessMultipleOrderCreatedMessages(t *testing.T) {
	truncateNotifications(testDB)

	// ===== create multiple messages =====
	messages := []dto.OrderCreatedEvent{
		{
			OrderID:     uuid.Must(uuid.NewV4()),
			Email:       "user1@example.com",
			TotalAmount: 100000,
		},
		{
			OrderID:     uuid.Must(uuid.NewV4()),
			Email:       "user2@example.com",
			TotalAmount: 200000,
		},
		{
			OrderID:     uuid.Must(uuid.NewV4()),
			Email:       "user3@example.com",
			TotalAmount: 300000,
		},
	}

	ctx := context.Background()

	// ===== process all messages =====
	for _, msg := range messages {
		messageBytes, err := json.Marshal(msg)
		assert.Nil(t, err)

		err = testController.Create(ctx, messageBytes)
		assert.Nil(t, err)
	}

	// ===== assert all records in database =====
	var count int64
	testDB.Model(&model.Notifications{}).Count(&count)
	assert.Equal(t, int64(3), count)

	// ===== verify each notification =====
	for _, msg := range messages {
		var notification model.Notifications
		err := testDB.First(&notification, "order_id = ?", msg.OrderID).Error
		assert.Nil(t, err)
		assert.Equal(t, "order_created", notification.Type)
	}
}

func TestProcessOrderCreatedMessageWithDifferentFormats(t *testing.T) {
	truncateNotifications(testDB)

	testCases := []struct {
		name          string
		message       string
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid PascalCase JSON",
			message:     `{"OrderID":"` + uuid.Must(uuid.NewV4()).String() + `","Email":"test@example.com","TotalAmount":100000}`,
			expectError: false,
		},
		{
			name:          "Missing OrderID field",
			message:       `{"Email":"test@example.com","TotalAmount":100000}`,
			expectError:   true,
			errorContains: "order_id is required",
		},
		{
			name:          "Missing Email field",
			message:       `{"OrderID":"` + uuid.Must(uuid.NewV4()).String() + `","TotalAmount":100000}`,
			expectError:   false, // Email tidak wajib di event
		},
		{
			name:        "With zero TotalAmount",
			message:     `{"OrderID":"` + uuid.Must(uuid.NewV4()).String() + `","Email":"test@example.com","TotalAmount":0}`,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			err := testController.Create(ctx, []byte(tc.message))

			if tc.expectError {
				assert.NotNil(t, err)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestProcessOrderCreatedMessageConcurrently(t *testing.T) {
	truncateNotifications(testDB)

	// ===== create messages for concurrent processing =====
	messageCount := 10
	messages := make([][]byte, messageCount)

	for i := 0; i < messageCount; i++ {
		msg := dto.OrderCreatedEvent{
			OrderID:     uuid.Must(uuid.NewV4()),
			Email:       "concurrent@example.com",
			TotalAmount: float64((i + 1) * 100000),
		}
		messageBytes, err := json.Marshal(msg)
		assert.Nil(t, err)
		messages[i] = messageBytes
	}

	// ===== process messages concurrently =====
	ctx := context.Background()
	done := make(chan bool, messageCount)

	for _, msg := range messages {
		go func(message []byte) {
			err := testController.Create(ctx, message)
			assert.Nil(t, err)
			done <- true
		}(msg)
	}

	// ===== wait for all goroutines =====
	for i := 0; i < messageCount; i++ {
		<-done
	}

	// ===== assert all records created =====
	var count int64
	testDB.Model(&model.Notifications{}).Count(&count)
	assert.Equal(t, int64(messageCount), count)
}

// ========== HELPER TEST FOR MESSAGE FORMAT ==========

func TestOrderCreatedEventJSONMarshaling(t *testing.T) {
	event := dto.OrderCreatedEvent{
		OrderID:     uuid.Must(uuid.NewV4()),
		Email:       "test@example.com",
		TotalAmount: 500000,
	}

	// ===== marshal to JSON =====
	jsonBytes, err := json.Marshal(event)
	assert.Nil(t, err)

	// ===== unmarshal back =====
	var decoded dto.OrderCreatedEvent
	err = json.Unmarshal(jsonBytes, &decoded)
	assert.Nil(t, err)

	// ===== assert values =====
	assert.Equal(t, event.OrderID, decoded.OrderID)
	assert.Equal(t, event.Email, decoded.Email)
	assert.Equal(t, event.TotalAmount, decoded.TotalAmount)
}

func TestOrderCreatedEventWithRealRabbitMQFormat(t *testing.T) {
	truncateNotifications(testDB)

	// ===== simulate exact format from RabbitMQ =====
	realMessage := `{"OrderID":"127fc70b-f1b7-465c-9cf4-9ee768fc8a56","Email":"dani@mail.com","TotalAmount":300000}`

	// ===== process message =====
	ctx := context.Background()
	err := testController.Create(ctx, []byte(realMessage))

	// ===== assert success =====
	assert.Nil(t, err)

	// ===== verify in database =====
	var notification model.Notifications
	err = testDB.First(&notification, "order_id = ?", "127fc70b-f1b7-465c-9cf4-9ee768fc8a56").Error
	assert.Nil(t, err)

	assert.Equal(t, "order_created", notification.Type)
	assert.Equal(t, "Your order has been created.", notification.Message)
}

// ========== TEST HTTP ENDPOINTS ==========

func TestGetAllNotificationsSuccess(t *testing.T) {
	truncateNotifications(testDB)

	// ===== create test data via GORM =====
	notifications := []model.Notifications{
		{
			OrderID: uuid.Must(uuid.NewV4()),
			Type:    "order_created",
			Message: "Your order has been created.",
			IsRead:  false,
		},
		{
			OrderID: uuid.Must(uuid.NewV4()),
			Type:    "order_created",
			Message: "Your order has been created.",
			IsRead:  true,
		},
		{
			OrderID: uuid.Must(uuid.NewV4()),
			Type:    "order_created",
			Message: "Your order has been created.",
			IsRead:  false,
		},
	}

	for _, notif := range notifications {
		err := testDB.Create(&notif).Error
		assert.Nil(t, err)
	}

	// ===== setup router =====
	router := gin.Default()
	router.GET("/api/v1/notifications", testController.GetAll)

	// ===== request GET =====
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/notifications",
		nil,
	)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// ===== assert response =====
	resp := recorder.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)

	var response map[string]any
	_ = json.Unmarshal(respBody, &response)

	assert.Equal(t, "OK", response["status"])

	data := response["data"].([]any)
	assert.Len(t, data, 3)

	// ===== verify first notification =====
	item := data[0].(map[string]any)
	assert.NotNil(t, item["id"])
	assert.NotNil(t, item["order_id"])
	assert.Equal(t, "order_created", item["type"])
	assert.Equal(t, "Your order has been created.", item["message"])
	assert.NotNil(t, item["is_read"])
	assert.NotNil(t, item["created_at"])
}

func TestGetAllNotificationsEmpty(t *testing.T) {
	truncateNotifications(testDB)

	// ===== setup router =====
	router := gin.Default()
	router.GET("/api/v1/notifications", testController.GetAll)

	// ===== request GET =====
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/notifications",
		nil,
	)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// ===== assert response =====
	resp := recorder.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)

	var response map[string]any
	_ = json.Unmarshal(respBody, &response)

	assert.Equal(t, "OK", response["status"])

	data := response["data"].([]any)
	assert.Len(t, data, 0)
}

func TestGetAllNotificationsFilterByIsRead(t *testing.T) {
	truncateNotifications(testDB)

	// ===== create mixed read/unread notifications =====
	notifications := []model.Notifications{
		{
			OrderID: uuid.Must(uuid.NewV4()),
			Type:    "order_created",
			Message: "Your order has been created.",
			IsRead:  false,
		},
		{
			OrderID: uuid.Must(uuid.NewV4()),
			Type:    "order_created",
			Message: "Your order has been created.",
			IsRead:  false,
		},
		{
			OrderID: uuid.Must(uuid.NewV4()),
			Type:    "order_created",
			Message: "Your order has been created.",
			IsRead:  true,
		},
	}

	for _, notif := range notifications {
		err := testDB.Create(&notif).Error
		assert.Nil(t, err)
	}

	// ===== setup router =====
	router := gin.Default()
	router.GET("/api/v1/notifications", testController.GetAll)

	// ===== request GET with query parameter =====
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/notifications?is_read=false",
		nil,
	)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// ===== assert response =====
	resp := recorder.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)

	var response map[string]any
	_ = json.Unmarshal(respBody, &response)

	assert.Equal(t, "OK", response["status"])

	data := response["data"].([]any)
	assert.Len(t, data, 2)

	// ===== verify all are unread =====
	for _, d := range data {
		item := d.(map[string]any)
		assert.False(t, item["is_read"].(bool))
	}
}

func TestGetAllNotificationsPagination(t *testing.T) {
	truncateNotifications(testDB)

	// ===== create multiple notifications =====
	for i := 0; i < 15; i++ {
		notif := model.Notifications{
			OrderID: uuid.Must(uuid.NewV4()),
			Type:    "order_created",
			Message: "Your order has been created.",
			IsRead:  false,
		}
		err := testDB.Create(&notif).Error
		assert.Nil(t, err)
	}

	// ===== setup router =====
	router := gin.Default()
	router.GET("/api/v1/notifications", testController.GetAll)

	// ===== request with pagination =====
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/notifications?page=1&limit=10",
		nil,
	)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// ===== assert response =====
	resp := recorder.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)

	var response map[string]any
	_ = json.Unmarshal(respBody, &response)

	assert.Equal(t, "OK", response["status"])

	data := response["data"].([]any)
	assert.Len(t, data, 10)
}