package notification

import "context"

type MessageSender interface {
	Send(ctx context.Context, message string) error
}