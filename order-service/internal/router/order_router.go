package router

import (
	"order_service/internal/controller"
	"order_service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(controller controller.OrderController) *gin.Engine {
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