package main

import (
	"context"
	"github.com/water2027/webhook/internal/application"
	"github.com/water2027/webhook/internal/infrastructure/config"
	"github.com/water2027/webhook/internal/infrastructure/feishu"
	"github.com/water2027/webhook/internal/infrastructure/persistence"
	"log"
	"github.com/water2027/webhook/internal/infrastructure/http_server"
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

	// 6. 初始化应用层编排
	webhookApp := application.NewWebhookApp(sourceRepo, messageSender)

	port := config.GlobalConfig.Port
	server := http_server.NewHTTPServer(webhookApp)
	
	log.Printf("Server starting on :%s...\n", port)
	if err := server.Start(config.GlobalConfig.Port); err != nil {
		log.Fatal(err)
	}
}
