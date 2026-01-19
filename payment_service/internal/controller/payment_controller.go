package controller

import (
	"context"

	"github.com/gin-gonic/gin"
)

type PaymentController interface {
	Create(ctx context.Context, body []byte) error
	HandleMidtransWebhook(ctx *gin.Context)
}