package feishu

import (
	"context"
	"fmt"

	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/water2027/webhook/internal/domain/notification"
	"github.com/water2027/webhook/internal/infrastructure/config"
	"github.com/water2027/webhook/internal/interfaces"
)

type larkBot struct {
	receiveId string
	client    *lark.Client
}

func NewLarkBot() interfaces.MessageSender {
	client := lark.NewClient(config.GlobalConfig.FeishuAppID, config.GlobalConfig.FeishuAppSecret)
	receiveId := config.GlobalConfig.FeishuOpenID
	return &larkBot{client: client, receiveId: receiveId}
}

func (b *larkBot) Send(ctx context.Context, message *notification.Message) error {
	msg := fmt.Sprintf("来自: %s\n标题: %s\n内容: %s\n时间: %s", message.Source, message.Title, message.Content, message.Date)
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(`open_id`).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(b.receiveId).
			MsgType(`text`).
			Content(fmt.Sprintf(`{"text":"%s"}`, msg)).
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
