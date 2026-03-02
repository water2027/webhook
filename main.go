package main

import (
	"context"
	"encoding/json"
	"github.com/water2027/webhook/internal/application"
	"github.com/water2027/webhook/internal/domain/webhook"
	"github.com/water2027/webhook/internal/infrastructure/config"
	"github.com/water2027/webhook/internal/infrastructure/feishu"
	"github.com/water2027/webhook/internal/infrastructure/persistence"
	"log"
	"net/http"
)

func main() {
	// 1. 初始化配置
	config.Load()

	// 2. 初始化数据库连接池
	ctx := context.Background()
	dbPool, err := persistence.NewPostgresPool(ctx, config.GlobalConfig.DatabaseURL)
	if err != nil {
		log.Fatalf("Database initialization failed: %v\n", err)
	}
	defer dbPool.Close()

	// 3. 初始化数据库 Schema
	if err := persistence.InitSchema(ctx, dbPool); err != nil {
		log.Fatalf("Schema initialization failed: %v\n", err)
	}

	// 4. 初始化基础设施
	sourceRepo := persistence.NewPostgresSourceRepository(dbPool)
	messageSender := feishu.NewLarkBot()

	// 5. 初始化领域服务
	webhookService := webhook.NewService()

	// 6. 初始化应用层编排
	webhookApp := application.NewWebhookApp(sourceRepo, webhookService, messageSender)

	// 7. 注册路由
	http.HandleFunc("/sources", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var body struct {
				Name string `json:"name"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			src, err := webhookApp.RegisterSource(r.Context(), body.Name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(src)
		}
	})

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var payload webhook.Payload
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			timestamp := r.Header.Get("X-Webhook-Timestamp")
			signature := r.Header.Get("X-Webhook-Signature")

			if err := webhookApp.Handle(r.Context(), &payload, timestamp, signature); err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	})

	port := config.GlobalConfig.Port
	log.Printf("Server starting on :%s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
