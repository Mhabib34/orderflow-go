//go:build wireinject
// +build wireinject

package wire

import (
	"order_service/internal/broker"
	"order_service/internal/controller"
	"order_service/internal/database"
	"order_service/internal/repository"
	"order_service/internal/router"
	"order_service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Router *gin.Engine
}

func NewValidator() *validator.Validate {
	return validator.New()
}

func NewRabbitMQ() *broker.RabbitMQ {
	return broker.NewRabbitMQ("amqp://guest:guest@localhost:5672/")
}

var repositorySet = wire.NewSet(
	repository.NewOrderRepository,
)

var usecaseSet = wire.NewSet(
	usecase.NewOrderUsecase,
)

var controllerSet = wire.NewSet(
	controller.NewOrderController,
)

var routerSet = wire.NewSet(
	router.SetupRouter,
)

func InitializeServer() (*App, error) {
	wire.Build(
		// Database
		database.Connect,

		// RabbitMQ
		NewRabbitMQ,

		// Validator
		NewValidator,

		// Layers
		repositorySet,
		usecaseSet,
		controllerSet,
		routerSet,

		// App
		wire.Struct(new(App), "*"),
	)
	return nil, nil
}
