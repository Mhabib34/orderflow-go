package broker

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
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