package interfaces

import (
	"context"

	"github.com/water2027/webhook/internal/domain/notification"
	"github.com/water2027/webhook/internal/domain/source"
)

type WebhookApp interface {
	Handle(ctx context.Context, payload *notification.Message, secret string) error
	RegisterSource(ctx context.Context, name string) (*source.Source, error)
}

type SourceRepository interface {
	Save(ctx context.Context, s *source.Source) error
	FindByID(ctx context.Context, id string) (*source.Source, error)
}

type MessageSender interface {
	Send(ctx context.Context, message *notification.Message) error
}