// Package memory provides in-memory adapters for tests and local runs.
package memory

import (
	"context"
	"sync"

	"github.com/7cout/go-template/internal/domain"
)

// Items is a concurrency-safe in-memory item repository.
type Items struct {
	mu     sync.RWMutex
	nextID int64
	items  map[int64]domain.Item
}

// NewItems returns an empty repository.
func NewItems() *Items {
	return &Items{items: make(map[int64]domain.Item)}
}

// Create assigns an ID and stores the item.
func (r *Items) Create(_ context.Context, item domain.Item) (domain.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	item.ID = r.nextID
	r.items[item.ID] = item
	return item, nil
}

// Get returns domain.ErrNotFound for unknown IDs.
func (r *Items) Get(_ context.Context, id int64) (domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[id]
	if !ok {
		return domain.Item{}, domain.ErrNotFound
	}
	return item, nil
}
