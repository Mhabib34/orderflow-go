package usecase

import (
	"context"
	"encoding/json"
	"log"
	"order_service/internal/broker"
	"order_service/internal/dto"
	"order_service/internal/exception"
	"order_service/internal/helper"
	"order_service/internal/model"
	"order_service/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
)

type OrderUsecaseImpl struct {
	OrderRepository repository.OrderRepository
	Validate       *validator.Validate
	Publisher      broker.Publisher

}

func NewOrderUsecase(orderRepository repository.OrderRepository, validate *validator.Validate, publisher broker.Publisher) OrderUsecase {
	return &OrderUsecaseImpl{OrderRepository: orderRepository, Validate: validate, Publisher: publisher}
}

func (service *OrderUsecaseImpl) CreateOrder(ctx context.Context, request dto.CreateOrderRequest) (dto.OrderResponse, error){
	err := service.Validate.Struct(request)
	exception.PanicIfError(err)

	order := &model.Orders{TotalAmount: request.TotalAmount, Email: request.Email}

	order, err = service.OrderRepository.CreateOrder(ctx, order)
	exception.PanicIfError(err)
	
	//publish to rabbitmq
	event := OrderCreatedEvent{OrderID: order.ID.String(), Email: request.Email, TotalAmount: order.TotalAmount}
	body, err := json.Marshal(event)
	exception.PanicIfError(err)

	// publish
	err = service.Publisher.Publish(ctx, "order.created", body)
	exception.PanicIfError(err) 

	return helper.ToOrderResponse(*order), nil
}

func (service *OrderUsecaseImpl) FindByID(ctx context.Context, id uuid.UUID) (dto.OrderResponse, error) {
	order, err := service.OrderRepository.FindByID(ctx, id)
	exception.PanicIfError(err)

	return helper.ToOrderResponse(*order), nil
}

func (service *OrderUsecaseImpl) GetAll(ctx context.Context, status string, limit, page int) ([]model.Orders, int64, error) {
	orders, total, err := service.OrderRepository.GetAll(ctx, status, limit, page)
	exception.PanicIfError(err)

	return orders, total, nil
}

func (service *OrderUsecaseImpl) UpdateStatus(ctx context.Context, id uuid.UUID, request dto.UpdateStatusRequest) (dto.OrderResponse, error) {
	err := service.Validate.Struct(request)
	exception.PanicIfError(err)

	order, err := service.OrderRepository.FindByID(ctx, id)
	exception.PanicIfError(err)

	order.Status = string(request.Status)

	order, err = service.OrderRepository.UpdateStatus(ctx, id, order)
	exception.PanicIfError(err)

	return helper.ToOrderResponse(*order), nil
}
func (u *OrderUsecaseImpl) UpdateOrderStatus(ctx context.Context, event dto.PaymentStatusChangedEvent) error {
	// Map payment status ke order status
	var orderStatus string
	switch event.PaymentStatus {
	case "SUCCESS":
		orderStatus = "SUCCESS"
	case "FAILED":
		orderStatus = "CANCELLED"
	case "EXPIRED":
		orderStatus = "EXPIRED"
	default:
		orderStatus = "PENDING"
	}

	// ✅ Buat object Orders dengan semua field yang mau diupdate
	orderUpdate := &model.Orders{
		Status:        orderStatus,
		PaymentID:     &event.PaymentID,      // Update payment_id
		PaymentMethod: event.PaymentMethod, 
		PaymentStatus: orderStatus,   // Update payment_status
	}

	// Update order di database
	_, err := u.OrderRepository.UpdateStatus(ctx, event.OrderID, orderUpdate)
	if err != nil {
		log.Printf("❌ Failed to update order: %v\n", err)
		return err
	}

	log.Printf("✅ Order %s updated:\n", event.OrderID)
	log.Printf("   Status: %s\n", orderStatus)
	log.Printf("   PaymentID: %s\n", event.PaymentID)
	log.Printf("   PaymentMethod: %s\n", event.PaymentMethod)
	
	return nil
}