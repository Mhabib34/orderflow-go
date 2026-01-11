package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	consumer "payment_service/internal/broker"
	"payment_service/internal/controller"
	"syscall"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ===== Graceful shutdown =====
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Println("🛑 Shutdown signal received")
		cancel()
	}()

	// ===== Init controller =====
	paymentController := controller.NewPaymentController()

	// ===== Init RabbitMQ Consumer =====
	rabbitURL := "amqp://guest:guest@localhost:5672/"
	consumer, err := consumer.NewConsumer(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}

	// ===== Start consumer =====
	err = consumer.Start(
		ctx,
		5, // worker count
		paymentController.Create,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Payment Service is running")

	// ===== block main =====
	<-ctx.Done()
	log.Println("👋 Payment Service stopped")
}