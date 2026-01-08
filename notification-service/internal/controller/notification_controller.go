package controller

import (
	"context"

	"github.com/gin-gonic/gin"
)

type NotificationController interface {
	Create(ctx context.Context, body []byte) error
	GetAll(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	Update(ctx *gin.Context)
}