package source

import (
	"context"
)

type SourceRepository interface {
	Save(ctx context.Context, s *Source) error
	FindByID(ctx context.Context, id string) (*Source, error)
}
