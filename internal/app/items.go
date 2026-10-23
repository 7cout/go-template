// Package app implements use cases. It depends on domain and declares the
// interfaces it needs from adapters.
package app

import (
	"context"
	"fmt"

	"github.com/7cout/go-template/internal/domain"
)

// ItemRepository is declared here, at the consumer, and kept small.
type ItemRepository interface {
	Create(ctx context.Context, item domain.Item) (domain.Item, error)
	Get(ctx context.Context, id int64) (domain.Item, error)
}

// Items is the item use-case service.
type Items struct {
	repo ItemRepository
}

// NewItems returns a concrete service, not an interface.
func NewItems(repo ItemRepository) *Items {
	return &Items{repo: repo}
}

// Create validates the name and stores a new item.
func (s *Items) Create(ctx context.Context, name string) (domain.Item, error) {
	item, err := domain.NewItem(name)
	if err != nil {
		return domain.Item{}, fmt.Errorf("create item: %w", err)
	}
	created, err := s.repo.Create(ctx, item)
	if err != nil {
		return domain.Item{}, fmt.Errorf("create item: %w", err)
	}
	return created, nil
}

// Get returns an item by ID. A missing item yields domain.ErrNotFound.
func (s *Items) Get(ctx context.Context, id int64) (domain.Item, error) {
	item, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Item{}, fmt.Errorf("get item %d: %w", id, err)
	}
	return item, nil
}
