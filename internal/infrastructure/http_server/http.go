package http_server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/water2027/webhook/internal/domain/notification"
	"github.com/water2027/webhook/internal/interfaces"
)

type HTTPServer interface {
	Start(port string) error
}

type httpServer struct {
	webhookApp interfaces.WebhookApp
}

func NewHTTPServer(webhookApp interfaces.WebhookApp) *httpServer {
	return &httpServer{webhookApp: webhookApp}
}

func (s *httpServer) Start(port string) error {
	http.HandleFunc("/sources", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var body struct {
				Name string `json:"name"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			src, err := s.webhookApp.RegisterSource(r.Context(), body.Name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(src)
		}
	})

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var payload notification.Message
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			secret := r.Header.Get("X-Webhook-Secret")

			if err := s.webhookApp.Handle(r.Context(), &payload, secret); err != nil {
				fmt.Println(err.Error())
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	})

	return http.ListenAndServe(":"+port, nil)
}
