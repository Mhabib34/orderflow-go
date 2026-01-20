//go:build wireinject
// +build wireinject

package wire

import (
	"payment_service/internal/broker"
	"payment_service/internal/controller"
	"payment_service/internal/database"
	"payment_service/internal/repository"
	"payment_service/internal/router"
	"payment_service/internal/service"
	"payment_service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type App struct {
	DB       *gorm.DB
	Router   *gin.Engine
	Consumer *broker.RabbitMQ
	Controller controller.PaymentController
}

func NewValidator() *validator.Validate {
	return validator.New()
}

func NewRabbitMQ() (*broker.RabbitMQ, error) {
	return broker.NewConsumer("amqp://guest:guest@localhost:5672/")
}

func NewPublisher(rmq *broker.RabbitMQ) broker.Publisher {
	return rmq
}

var brokerSet = wire.NewSet(
	NewRabbitMQ,
	NewPublisher,
)

var repositorySet = wire.NewSet(
	repository.NewPaymentRepository,
)

var midtransServiceSet = wire.NewSet(
	service.NewMidtransService,
)

var usecaseSet = wire.NewSet(
	usecase.NewPaymentUsecase,
)

var controllerSet = wire.NewSet(
	controller.NewPaymentController,
)

var routerSet = wire.NewSet(
	router.SetupRouter,
)

func InitializeServer() (*App, error) {
	wire.Build(
		database.Connect,
		NewValidator,
		brokerSet,
		repositorySet,
		usecaseSet,
		midtransServiceSet,
		controllerSet,
		routerSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil
}