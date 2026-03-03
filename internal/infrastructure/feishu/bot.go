package feishu

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

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

type LarkTextContent struct {
	Text string `json:"text"`
}

func NewLarkBot() interfaces.MessageSender {
	client := lark.NewClient(config.GlobalConfig.FeishuAppID, config.GlobalConfig.FeishuAppSecret)
	receiveId := config.GlobalConfig.FeishuOpenID
	return &larkBot{client: client, receiveId: receiveId}
}

func (b *larkBot) Send(ctx context.Context, message *notification.Message) error {
	msg := fmt.Sprintf("来自: %s\n标题: %s\n内容: %s\n时间: %s", message.Source, message.Title, message.Content, message.Date)
	contentBytes, err := json.Marshal(LarkTextContent{Text: msg})
	if err != nil {
		slog.Error("Failed to marshal lark content", "error", err)
		return fmt.Errorf("failed to marshal lark content: %w", err)
	}
	contentStr := string(contentBytes)
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(`open_id`).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(b.receiveId).
			MsgType(`text`).
			Content(contentStr).
			Build()).
		Build()

	resp, err := b.client.Im.V1.Message.Create(ctx, req)
	if err != nil {
		slog.Error("Lark API call failed", "error", err)
		return err
	}

	if !resp.Success() {
		slog.Error("Lark message creation failed",
			"code", resp.Code,
			"msg", resp.Msg,
			"requestId", resp.RequestId(),
			"errorResponse", larkcore.Prettify(resp.CodeError),
		)
		return fmt.Errorf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
	}

	return nil
}
