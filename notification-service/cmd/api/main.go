package main

import (
	"context"
	"log"
	"notification_service/internal/consumer"
	"notification_service/internal/controller"
	"notification_service/internal/database"
	"notification_service/internal/repository"
	"notification_service/internal/router"
	"notification_service/internal/usecase"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env hanya sekali di awal
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}


func main() {
	rabbitURL := "amqp://guest:guest@localhost:5672/"
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	
	validate := validator.New()
	
	repo := repository.NewNotificationRepository(db)
	usecase := usecase.NewNotificationUsecase(repo, validate)
	controller := controller.NewNotificationController(usecase)
	router := router.SetupRouter(controller)

	consumer, err := consumer.NewConsumer(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}
	
	err = consumer.Start(context.Background(), 5, controller.Create)
	if err != nil {
		log.Fatal(err)
	} 

	log.Println("📥 Notification service started...")

	err = router.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("🛑 Shutting down consumer...")
}