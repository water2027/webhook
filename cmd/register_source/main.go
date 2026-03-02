package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/water2027/webhook/internal/domain/source"
	"github.com/water2027/webhook/internal/infrastructure/config"
	"github.com/water2027/webhook/internal/infrastructure/persistence"
)

func main() {
	name := flag.String("name", "", "Name of the source")
	flag.Parse()

	if *name == "" {
		fmt.Println("Usage: go run cmd/register_source/main.go -name <source_name>")
		os.Exit(1)
	}

	config.Load()

	ctx := context.Background()
	dbPool, err := persistence.NewPostgresPool(ctx, config.GlobalConfig.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	if err := persistence.InitSchema(ctx, dbPool); err != nil {
		log.Fatalf("Schema initialization failed: %v\n", err)
	}

	repo := persistence.NewPostgresSourceRepository(dbPool)

	s := source.NewSource(*name)
	if err := repo.Save(ctx, s); err != nil {
		log.Fatalf("Failed to save source: %v", err)
	}

	fmt.Println("Source registered successfully!")
	fmt.Printf("ID: %s\n", s.ID)
	fmt.Printf("Name: %s\n", s.Name)
	fmt.Printf("Secret: %s\n", s.Secret)
}
