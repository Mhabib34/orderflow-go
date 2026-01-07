package controller

import "context"

type NotificationController interface {
	Create(ctx context.Context, body []byte) error
}