package controller

import (
	"math"
	"net/http"
	"order_service/internal/dto"
	"order_service/internal/exception"
	"order_service/internal/helper"
	"order_service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type OrderControllerImpl struct {
	OrderUsecase usecase.OrderUsecase
}

func NewOrderController(orderUsecase usecase.OrderUsecase) OrderController {
	return &OrderControllerImpl{OrderUsecase: orderUsecase}
}

func (controller *OrderControllerImpl) CreateOrder(ctx *gin.Context) {
	var request dto.CreateOrderRequest

	if err := ctx.ShouldBind(&request); err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	result, err := controller.OrderUsecase.CreateOrder(ctx.Request.Context(), request)
	if err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	webResponse := dto.WebResponse{
		Status: "OK",
		Message: "Order created successfully", 
		Data: result,
	}
	

	helper.WriteToResponseBody(ctx,http.StatusCreated,webResponse)
}

func (controller *OrderControllerImpl) FindByID(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := helper.StringToUUID(idParam)
	if err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	result, err := controller.OrderUsecase.FindByID(ctx.Request.Context(), id)
	if err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	webResponse := dto.WebResponse{
		Status: "OK",
		Message: "Order found successfully", 
		Data: result,
	}

	helper.WriteToResponseBody(ctx, http.StatusOK, webResponse)
}

func (controller *OrderControllerImpl) GetAll(ctx *gin.Context) {
	page := helper.StringToIntDefault(ctx.Query("page"), 1)
	limit := helper.StringToIntDefault(ctx.Query("limit"), 10)
	status := ctx.Query("status")

	orders, total, err := controller.OrderUsecase.GetAll(
		ctx.Request.Context(), 
		status, 
		limit, 
		page,
	)
	if err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	webResponse := dto.WebResponse{
		Status: "OK",
		Message: "Orders fetched successfully", 
		Data: orders,
		Pagination: &dto.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}

	helper.WriteToResponseBody(ctx, http.StatusOK, webResponse)
}

func(controller *OrderControllerImpl) UpdateStatus(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := helper.StringToUUID(idParam)
	if err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	var request dto.UpdateStatusRequest

	if err := ctx.ShouldBind(&request); err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	result, err := controller.OrderUsecase.UpdateStatus(ctx.Request.Context(), id, request)
	if err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	webResponse := dto.WebResponse{
		Status: "OK",
		Message: "Order updated successfully", 
		Data: result,
	}

	helper.WriteToResponseBody(ctx, http.StatusOK, webResponse)
}
