package helper

import (
	"payment_service/internal/dto"
	"payment_service/internal/model"
)

func ToPaymentResponse(payment model.Payments) dto.PaymentResponse {
	return dto.PaymentResponse{
		ID:         payment.ID,
		OrderID:    payment.OrderID,
		Amount:     int64(payment.Amount),
		Method:     payment.Method,
		Status:     string(payment.Status),
		PaidAt:     payment.PaidAt.String(),
		ExpiredAt:  payment.ExpiredAt.String(),
		ProviderRefId: payment.ProviderRefID,
		CreatedAt:  payment.CreatedAt.String(),
		UpdatedAt:  payment.UpdatedAt.String(),
		PaymentURL: payment.PaymentURL,
		PaymentID:  payment.PaymentID,
	}
}