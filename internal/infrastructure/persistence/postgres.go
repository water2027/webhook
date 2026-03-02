package persistence

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func NewPostgresPool(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := dbPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return dbPool, nil
}

func InitSchema(ctx context.Context, dbPool *pgxpool.Pool) error {
	schema := `
	CREATE TABLE IF NOT EXISTS sources (
		id VARCHAR(16) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		secret VARCHAR(32) NOT NULL
	);
	`
	_, err := dbPool.Exec(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	log.Println("Database schema initialized successfully")
	return nil
}
