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

func NewRabbitMQ(url string) *RabbitMQ {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal("failed connect rabbitmq:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("failed open channel:", err)
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
		log.Fatal("failed declare exchange:", err)
	}

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}
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

	// ✅ MESSAGE SUCCESS
	fmt.Printf("📤 Message published | routing_key=%s\n", routingKey)

	return nil
}