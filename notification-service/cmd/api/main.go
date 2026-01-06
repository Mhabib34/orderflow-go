package main

import (
	"context"
	"log"
	"notification_service/internal/consumer"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	rabbitURL := "amqp://guest:guest@localhost:5672/"

	consumer, err := consumer.NewConsumer(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}
	
	err = consumer.Start(context.Background(), 5)
	if err != nil {
		log.Fatal(err)
	} 

	log.Println("📥 Notification service started...")

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("🛑 Shutting down consumer...")
}