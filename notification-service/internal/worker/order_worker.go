package worker

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Worker struct {
	ID      int
}

func (w *Worker) Start(
	ctx context.Context,
	jobs <-chan amqp.Delivery,
) {
	log.Printf("👷 Worker %d started\n", w.ID)
 
	for msg := range jobs {
		log.Printf("📩 Worker %d received message\n", w.ID)

		msg.Ack(false)
		log.Printf("✅ Worker %d message acked\n", w.ID)
	}
}
