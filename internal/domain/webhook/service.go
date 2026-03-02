package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Payload struct {
	EventID string `json:"eventId"`
	Source  string `json:"source"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Date    string `json:"date"`
}

type Service interface {
	Verify(payload *Payload, secret string, timestamp string, signature string) error
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) Verify(payload *Payload, secret string, timestamp string, signature string) error {
	if !verifySignature(secret, timestamp, payload, signature) {
		return fmt.Errorf("signature verification failed")
	}
	return nil
}

func verifySignature(secret string, timestamp string, payload *Payload, signature string) bool {
	h := hmac.New(sha256.New, []byte(secret))
	data := fmt.Sprintf("%s\n%s\n%s", timestamp, payload.Source, payload.EventID)
	h.Write([]byte(data))
	expectedSignature := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}
