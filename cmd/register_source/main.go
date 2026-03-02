package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
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

	config.InitConfig()

	dbURL := config.GetWithDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/webhook?sslmode=disable")
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	repo := persistence.NewPostgresSourceRepository(dbPool)

	s := source.NewSource(*name)
	if err := repo.Save(context.Background(), s); err != nil {
		log.Fatalf("Failed to save source: %v", err)
	}

	fmt.Println("Source registered successfully!")
	fmt.Printf("ID: %s\n", s.ID)
	fmt.Printf("Name: %s\n", s.Name)
	fmt.Printf("Secret: %s\n", s.Secret)
}
