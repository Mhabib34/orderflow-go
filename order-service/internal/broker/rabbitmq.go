package broker

import (
	"context"
	"fmt"
	"log"
	"order_service/internal/worker"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   string
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	log.Println("🔵 [DEBUG] Attempting to connect to RabbitMQ...")
	
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("🔴 [ERROR] Failed to connect to RabbitMQ: %v\n", err)
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Printf("🔴 [ERROR] Failed to open channel: %v\n", err)
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Queue
	q, err := ch.QueueDeclare(
		"order.payment.status.updated",
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
		"payment.status.updated",
		"payment_exchanges",
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

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
		ch.Close()
		conn.Close()
		log.Printf("🔴 [ERROR] Failed to declare exchange: %v\n", err)
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	log.Println("✅ [SUCCESS] RabbitMQ connected successfully")
	
	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (c *RabbitMQ) Start(
	ctx context.Context,
	workerCount int,
	controller func(context.Context, []byte) error,
	) error {

	err := c.Channel.Qos(workerCount, 0, false)
	if err != nil {
		return err
	}

	msgs, err := c.Channel.Consume(
		c.Queue,
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
			Controller: controller,
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

func (r *RabbitMQ) Publish(ctx context.Context, routingKey string, body []byte) error {
	err := r.Channel.PublishWithContext(
		ctx,
		"order_exchanges",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return err
	}

	fmt.Printf("📤 Message published | routing_key=%s\n", routingKey)
	return nil
}