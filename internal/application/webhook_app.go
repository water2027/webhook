package application

import (
	"context"
	"fmt"
	"github.com/water2027/webhook/internal/domain/notification"
	"github.com/water2027/webhook/internal/domain/source"
	"github.com/water2027/webhook/internal/domain/webhook"
)

type WebhookApp interface {
	Handle(ctx context.Context, payload *webhook.Payload, timestamp string, signature string) error
	RegisterSource(ctx context.Context, name string) (*source.Source, error)
}

type webhookApp struct {
	sourceRepo     source.SourceRepository
	webhookService webhook.Service
	notifier       notification.MessageSender
}

func NewWebhookApp(repo source.SourceRepository, ws webhook.Service, n notification.MessageSender) WebhookApp {
	return &webhookApp{
		sourceRepo:     repo,
		webhookService: ws,
		notifier:       n,
	}
}

func (a *webhookApp) Handle(ctx context.Context, payload *webhook.Payload, timestamp string, signature string) error {
	src, err := a.sourceRepo.FindByID(ctx, payload.Source)
	if err != nil {
		return fmt.Errorf("source not found: %w", err)
	}

	if err := a.webhookService.Verify(payload, src.Secret, timestamp, signature); err != nil {
		return err
	}

	msg := fmt.Sprintf("Source: %s\nTitle: %s\nContent: %s\nDate: %s", src.Name, payload.Title, payload.Content, payload.Date)
	return a.notifier.Send(ctx, msg)
}

func (a *webhookApp) RegisterSource(ctx context.Context, name string) (*source.Source, error) {
	s := source.NewSource(name)
	if err := a.sourceRepo.Save(ctx, s); err != nil {
		return nil, err
	}
	return s, nil
}
