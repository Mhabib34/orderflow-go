package helper

import (
	"order_service/internal/dto"
	"order_service/internal/model"
)

func ToOrderResponse(order model.Orders) dto.OrderResponse {
	return dto.OrderResponse{
		ID:          order.ID,
		Status:      order.Status,
		Email:       order.Email,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt,
	}
}