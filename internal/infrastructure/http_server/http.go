package http_server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

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
	mux := http.NewServeMux()

	mux.HandleFunc("/sources", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var body struct {
				Name string `json:"name"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				slog.Warn("Failed to decode source registration body", "error", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			src, err := s.webhookApp.RegisterSource(r.Context(), body.Name)
			if err != nil {
				slog.Error("Failed to register source", "name", body.Name, "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(src)
		}
	})

	mux.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var payload notification.Message
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				slog.Warn("Failed to decode webhook payload", "error", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			secret := r.Header.Get("X-Webhook-Secret")

			if err := s.webhookApp.Handle(r.Context(), &payload, secret); err != nil {
				slog.Error("Failed to handle webhook", "source", payload.Source, "title", payload.Title, "error", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	})

	// Wrap mux with logging middleware
	loggedMux := loggingMiddleware(mux)

	return http.ListenAndServe(":"+port, loggedMux)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("Request handled",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"duration", time.Since(start),
		)
	})
}
