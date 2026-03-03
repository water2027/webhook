package application

import (
	"context"
	"fmt"
	"github.com/water2027/webhook/internal/domain/notification"
	"github.com/water2027/webhook/internal/domain/source"
	"github.com/water2027/webhook/internal/interfaces"
)

type webhookApp struct {
	sourceRepo     interfaces.SourceRepository
	notifier       interfaces.MessageSender
}

func NewWebhookApp(repo interfaces.SourceRepository, n interfaces.MessageSender) interfaces.WebhookApp {
	return &webhookApp{
		sourceRepo:     repo,
		notifier:       n,
	}
}

func (a *webhookApp) Handle(ctx context.Context, payload *notification.Message, secret string) error {
	src, err := a.sourceRepo.FindByID(ctx, payload.Source)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("source not found: %w", err)
	}

	if !src.Verify(payload.Source, secret) {
		fmt.Println(payload.Source, src.ID)
		fmt.Println(secret, src.Secret)
		return fmt.Errorf("invalid source credentials")
	}

	return a.notifier.Send(ctx, payload)
}

func (a *webhookApp) RegisterSource(ctx context.Context, name string) (*source.Source, error) {
	s := source.NewSource(name)
	if err := a.sourceRepo.Save(ctx, s); err != nil {
		return nil, err
	}
	return s, nil
}
