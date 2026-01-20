package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"payment_service/cmd/wire"
	"payment_service/internal/config"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// Load midtrans config
	config.InitMidtrans()

	// Context untuk graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize dependencies dengan Wire
	app, err := wire.InitializeServer()
	if err != nil {
		log.Fatalf("❌ Failed to initialize server: %v", err)
	}

	// Start RabbitMQ Consumer untuk payment.status.updated
	err = app.Consumer.Start(ctx, 5, app.Controller.Create)
	if err != nil {
		log.Fatalf("❌ Failed to start RabbitMQ consumer: %v", err)
	}
	log.Println("✅ RabbitMQ consumer started")

	// HTTP Server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.Router,
	}

	// Start HTTP server di goroutine
	go func() {
		log.Println("✅ PaymentService is running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// Wait for interrupt signal (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutdown signal received, stopping server...")
	cancel() // Cancel context untuk stop RabbitMQ consumer

	// Graceful shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	}

	// Close RabbitMQ connections
	if app.Consumer.Channel != nil {
		app.Consumer.Channel.Close()
		log.Println("✅ RabbitMQ channel closed")
	}
	if app.Consumer.Conn != nil {
		app.Consumer.Conn.Close()
		log.Println("✅ RabbitMQ connection closed")
	}

	log.Println("👋 PaymentService stopped gracefully")
}