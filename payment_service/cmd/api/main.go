package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	consumer "payment_service/internal/broker"
	"payment_service/internal/config"
	"payment_service/internal/controller"
	"payment_service/internal/database"
	"payment_service/internal/repository"
	"payment_service/internal/router"
	"payment_service/internal/service"
	"payment_service/internal/usecase"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// 1. Database
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// 2. RabbitMQ - Inisialisasi LEBIH AWAL
	rabbitURL := "amqp://guest:guest@localhost:5672/"
	rabbitConsumer, err := consumer.NewConsumer(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Dependencies
	validator := validator.New()
	config.InitMidtrans()
	
	midtransService := service.NewMidtransService()
	repo := repository.NewPaymentRepository(db)
	
	// ✅ FIX: Pass rabbitConsumer sebagai Publisher
	paymentUsecase := usecase.NewPaymentUsecase(
		repo, 
		validator, 
		midtransService,
		rabbitConsumer, // ✅ Inject Publisher di sini
	)
	
	paymentController := controller.NewPaymentController(paymentUsecase)
	
	// 4. Router
	r := router.SetupRouter(*paymentController.(*controller.PaymentControllerImpl))
	
	// 5. Start RabbitMQ Consumer
	err = rabbitConsumer.Start(ctx, 5, paymentController.Create)
	if err != nil {
		log.Fatal(err)
	}

	// 6. HTTP Server dengan context
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// 7. Start server di goroutine
	go func() {
		log.Println("✅ Payment Service is running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v\n", err)
		}
	}()

	// 8. Wait for interrupt signal (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("🛑 Shutdown signal received, stopping server...")
	cancel() // Cancel context untuk stop RabbitMQ consumer

	// 9. Graceful shutdown HTTP server dengan timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v\n", err)
	}

	// 10. Close RabbitMQ connections
	if rabbitConsumer.Channel != nil {
		rabbitConsumer.Channel.Close()
	}
	if rabbitConsumer.Conn != nil {
		rabbitConsumer.Conn.Close()
	}

	log.Println("👋 Payment Service stopped gracefully")
}