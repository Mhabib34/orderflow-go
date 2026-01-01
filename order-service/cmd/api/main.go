package main

import (
	"log"
	"order_service/cmd/wire"

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
	
	app.Router.Run(":3000")
}