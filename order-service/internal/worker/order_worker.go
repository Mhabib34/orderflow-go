package worker

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Worker struct {
	ID         int
	Controller func(context.Context, []byte) error
}

func (w *Worker) Start(ctx context.Context, jobs <-chan amqp.Delivery) {
	for msg := range jobs {
		log.Printf("📩 Worker %d received message\n", w.ID)

		err := w.Controller(ctx, msg.Body)
		if err != nil {
			log.Printf("❌ Worker %d error: %v\n", w.ID, err)
			
			// CRITICAL: Jangan requeue jika error parsing/validation
			// Ini akan menghentikan infinite retry loop
			msg.Nack(false, false) // false = tidak requeue
			
			// Alternatif: bisa tambahkan logic untuk requeue dengan limit
			// if msg.Headers["x-retry-count"].(int) < 3 {
			//     msg.Nack(false, true) // requeue
			// } else {
			//     msg.Nack(false, false) // buang ke DLQ
			// }
		} else {
			log.Printf("✅ Worker %d processed successfully\n", w.ID)
			msg.Ack(false)
		}
	}
}