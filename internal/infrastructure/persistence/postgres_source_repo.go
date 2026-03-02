package persistence

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/water2027/webhook/internal/domain/source"
)

type postgresSourceRepository struct {
	db *pgxpool.Pool
}

func NewPostgresSourceRepository(db *pgxpool.Pool) source.SourceRepository {
	return &postgresSourceRepository{db: db}
}

func (r *postgresSourceRepository) Save(ctx context.Context, s *source.Source) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO sources (id, name, secret) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET name = $2, secret = $3",
		s.ID, s.Name, s.Secret,
	)
	return err
}

func (r *postgresSourceRepository) FindByID(ctx context.Context, id string) (*source.Source, error) {
	var s source.Source
	err := r.db.QueryRow(ctx,
		"SELECT id, name, secret FROM sources WHERE id = $1",
		id,
	).Scan(&s.ID, &s.Name, &s.Secret)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
