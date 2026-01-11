package main

import (
	"context"
	"log"
	"notification_service/cmd/wire"
	"os"
	"os/signal"
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
	app, err := wire.InitializeServer()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Consumer.Start(
		context.Background(),
		5,
		app.Controller.Create,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("📥 Notification service started...")

	go func() {
		if err := app.Router.Run(":8000"); err != nil {
			log.Fatal(err)
		}
	}()

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("🛑 Shutting down consumer...")
}