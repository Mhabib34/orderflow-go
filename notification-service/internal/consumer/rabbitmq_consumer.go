package consumer

import (
	"context"
	"log"

	"notification_service/internal/worker"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewConsumer(url string) (*Consumer, error) {
	log.Println("🔌 Connecting to RabbitMQ...")

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Exchange
	err = ch.ExchangeDeclare(
		"order_exchanges",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Queue
	q, err := ch.QueueDeclare(
		"notification.order.created",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Binding
	err = ch.QueueBind(
		q.Name,
		"order.created",
		"order_exchanges",
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	log.Println("✅ RabbitMQ consumer ready")

	return &Consumer{
		conn:    conn,
		channel: ch,
		queue:   q.Name,
	}, nil
}

func (c *Consumer) Start(
	ctx context.Context,
	workerCount int,
) error {

	err := c.channel.Qos(workerCount, 0, false)
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		c.queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	jobs := make(chan amqp.Delivery)

	// ===== start workers =====
	for i := 1; i <= workerCount; i++ {
		w := worker.Worker{
			ID:      i,
		}
		go w.Start(ctx, jobs)
	}

	// ===== dispatcher =====
	go func() {
		for msg := range msgs {
			jobs <- msg
		}
	}()

	log.Printf("🚀 Consumer running with %d workers\n", workerCount)

	return nil
}
