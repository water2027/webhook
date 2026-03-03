package persistence

import (
	"context"
	"errors"
	"github.com/water2027/webhook/internal/domain/source"
	"github.com/water2027/webhook/internal/interfaces"
	"sync"
)

type memorySourceRepository struct {
	mu      sync.RWMutex
	sources map[string]*source.Source
}

func NewMemorySourceRepository() interfaces.SourceRepository {
	return &memorySourceRepository{
		sources: make(map[string]*source.Source),
	}
}

func (r *memorySourceRepository) Save(ctx context.Context, s *source.Source) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sources[s.ID] = s
	return nil
}

func (r *memorySourceRepository) FindByID(ctx context.Context, id string) (*source.Source, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	src, ok := r.sources[id]
	if !ok {
		return nil, errors.New("source not found")
	}
	return src, nil
}
