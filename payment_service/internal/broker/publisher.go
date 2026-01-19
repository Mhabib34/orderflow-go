package broker

import "context"

type Publisher interface {
	Publish(ctx context.Context, routingKey string, body []byte) error
}
