package controller

import (
	"context"

	"github.com/gin-gonic/gin"
)

type OrderController interface {
	CreateOrder(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	UpdateStatus(ctx *gin.Context)
	HandlePaymentStatusUpdated(ctx context.Context, body []byte) error 
}