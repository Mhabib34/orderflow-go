package router

import (
	"notification_service/internal/controller"
	"notification_service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(controller controller.NotificationController) *gin.Engine {
	r := gin.New()

	// middleware
	r.Use(gin.Logger())
	r.Use(middleware.ErrorRecovery()) // ⬅️ penting

	api := r.Group("/api/v1")
	{
		api.GET("/notifications", controller.GetAll)
		api.GET("/notifications/:id", controller.FindByID)
		api.PATCH("/notifications/:id", controller.Update)
	}

	return r
}