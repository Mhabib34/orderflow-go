package test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"order_service/internal/controller"
	"order_service/internal/entity"
	"order_service/internal/middleware"
	"order_service/internal/model"
	"order_service/internal/repository"
	"order_service/internal/usecase"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	testDB     *gorm.DB
	testRouter http.Handler
)

type FakePublisher struct{}

func (f *FakePublisher) Publish(ctx context.Context, routingKey string, body []byte) error {
	return nil
}

func setupTestDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=habib123 dbname=order_service_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&model.Orders{})
	if err != nil {
		panic(err)
	}

	return db
}

func setupRouter(db *gorm.DB) http.Handler {
	validate := validator.New()
	publisher := &FakePublisher{}

	repo := repository.NewOrderRepository(db)
	usecase := usecase.NewOrderUsecase(repo, validate, publisher)
	controller := controller.NewOrderController(usecase)

	r := gin.New()

	// middleware
	r.Use(gin.Logger())
	r.Use(middleware.ErrorRecovery()) // ⬅️ penting

	api := r.Group("/api/v1")
	{
		api.POST("/orders", controller.CreateOrder)
	}

	return r
}

func truncateMissingPersons(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE orders")
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	testDB = setupTestDB()
	testRouter = setupRouter(testDB)

	code := m.Run()

	os.Exit(code)
}

func TestCreateOrderSuccess(t *testing.T) {
	truncateMissingPersons(testDB)

	// ===== multipart body =====
	payload := map[string]any{
	"email":        "QzZ0s@example.com",
	"total_amount": 100000,
}

	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/orders",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	// ===== assert response =====
	resp := recorder.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)

	var response map[string]any
	json.Unmarshal(respBody, &response)

	assert.Equal(t, "OK", response["status"])

	data := response["data"].(map[string]any)

	assert.Equal(t, "QzZ0s@example.com", data["email"])
	assert.Equal(t, "pending", data["status"])
	assert.Equal(t, float64(100000), data["total_amount"])


	// ===== assert DB =====
	var count int64
	testDB.Model(&entity.Orders{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestCreateOrderFailBadRequest(t *testing.T) {
	truncateMissingPersons(testDB)

	// ===== multipart body =====
	payload := map[string]any{
	"email":        "",
	"total_amount": 100000,
}

	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/orders",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	// ===== assert response =====
	resp := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)

	var response map[string]any
	json.Unmarshal(respBody, &response)

	assert.Equal(t, "BAD REQUEST", response["status"])
}