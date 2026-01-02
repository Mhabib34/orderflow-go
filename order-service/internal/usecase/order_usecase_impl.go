package usecase

import (
	"context"
	"encoding/json"
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