package router

import (
	"payment_service/internal/controller"
	"payment_service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(controller controller.PaymentControllerImpl) *gin.Engine {
	r := gin.New()

	// middleware
	r.Use(gin.Logger())
	r.Use(middleware.ErrorRecovery()) // ⬅️ penting

	api := r.Group("/api/v1")
	{
		api.POST("/payments/webhook", controller.HandleMidtransWebhook)
	}

	return r
}