package main

import (
	"context"
	"log/slog"

	"github.com/water2027/webhook/internal/application"
	"github.com/water2027/webhook/internal/infrastructure/config"
	"github.com/water2027/webhook/internal/infrastructure/feishu"
	"github.com/water2027/webhook/internal/infrastructure/http_server"
	"github.com/water2027/webhook/internal/infrastructure/logger"
	"github.com/water2027/webhook/internal/infrastructure/persistence"
)

func main() {
	// 1. 初始化配置
	config.Load()

	// 2. 初始化日志
	logger.Init()

	// 3. 初始化数据库连接池
	ctx := context.Background()
	dbPool, err := persistence.NewPostgresPool(ctx, config.GlobalConfig.DatabaseURL)
	if err != nil {
		slog.Error("Database initialization failed", "error", err)
		return
	}
	defer dbPool.Close()

	// 4. 初始化数据库 Schema
	if err := persistence.InitSchema(ctx, dbPool); err != nil {
		slog.Error("Schema initialization failed", "error", err)
		return
	}

	// 5. 初始化基础设施
	sourceRepo := persistence.NewPostgresSourceRepository(dbPool)
	messageSender := feishu.NewLarkBot()

	// 6. 初始化应用层编排
	webhookApp := application.NewWebhookApp(sourceRepo, messageSender)

	port := config.GlobalConfig.Port
	server := http_server.NewHTTPServer(webhookApp)

	slog.Info("Server starting", "port", port)
	if err := server.Start(config.GlobalConfig.Port); err != nil {
		slog.Error("Server failed", "error", err)
	}
}
