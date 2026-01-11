//go:build wireinject
// +build wireinject

package wire

import (
	"notification_service/internal/consumer"
	"notification_service/internal/controller"
	"notification_service/internal/database"
	"notification_service/internal/repository"
	"notification_service/internal/router"
	"notification_service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type App struct {
	DB         *gorm.DB
	Router     *gin.Engine
	Consumer   *consumer.Consumer
	Controller controller.NotificationController
}

/* ========= INFRA ========= */

func NewValidator() *validator.Validate {
	return validator.New()
}

func NewRabbitURL() string {
	return "amqp://guest:guest@localhost:5672/"
}

/* ========= SETS ========= */

var consumerSet = wire.NewSet(
	NewRabbitURL,
	consumer.NewConsumer,
)

var repositorySet = wire.NewSet(
	repository.NewNotificationRepository,
)

var usecaseSet = wire.NewSet(
	usecase.NewNotificationUsecase,
)

var controllerSet = wire.NewSet(
	controller.NewNotificationController,
)

var routerSet = wire.NewSet(
	router.SetupRouter,
)

/* ========= INJECTOR ========= */

func InitializeServer() (*App, error) {
	wire.Build(
		database.Connect,
		NewValidator,
		consumerSet,
		repositorySet,
		usecaseSet,
		controllerSet,
		routerSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil
}
