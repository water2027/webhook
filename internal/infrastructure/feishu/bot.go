package feishu

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/water2027/webhook/internal/domain/notification"
	"github.com/water2027/webhook/internal/infrastructure/config"
)

type larkBot struct {
	receiveId string
	client    *lark.Client
}

func NewLarkBot() notification.MessageSender {
	client := lark.NewClient(config.Get("FEISHU_APP_ID"), config.Get("FEISHU_APP_SECRET"))
	receiveId := config.Get("FEISHU_OPEN_ID")
	return &larkBot{client: client, receiveId: receiveId}
}

func (b *larkBot) Send(ctx context.Context, message string) error {
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(`open_id`).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(b.receiveId).
			MsgType(`text`).
			Content(fmt.Sprintf(`{"text":"%s"}`, message)).
			Build()).
		Build()

	resp, err := b.client.Im.V1.Message.Create(ctx, req)
	if err != nil {
		return err
	}

	if !resp.Success() {
		return fmt.Errorf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
	}

	return nil
}
