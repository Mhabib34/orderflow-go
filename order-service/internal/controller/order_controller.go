package controller

import "github.com/gin-gonic/gin"

type OrderController interface {
	CreateOrder(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	GetAll(ctx *gin.Context)
}