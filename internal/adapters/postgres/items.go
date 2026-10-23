// Package postgres implements repositories on PostgreSQL via sqlc-generated
// queries (internal/adapters/postgres/sqlcgen, regenerate with `task sqlc`).
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/7cout/go-template/internal/adapters/postgres/sqlcgen"
	"github.com/7cout/go-template/internal/domain"
)

// Items is a PostgreSQL item repository.
type Items struct {
	q *sqlcgen.Queries
}

// NewItems accepts a pool, a connection or a transaction.
func NewItems(db sqlcgen.DBTX) *Items {
	return &Items{q: sqlcgen.New(db)}
}

// Create inserts the item and returns it with the generated ID.
func (r *Items) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	row, err := r.q.CreateItem(ctx, item.Name)
	if err != nil {
		return domain.Item{}, fmt.Errorf("insert item: %v", err) // %v: do not expose driver errors
	}
	return domain.Item{ID: row.ID, Name: row.Name}, nil
}

// Get maps pgx.ErrNoRows to domain.ErrNotFound.
func (r *Items) Get(ctx context.Context, id int64) (domain.Item, error) {
	row, err := r.q.GetItem(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Item{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Item{}, fmt.Errorf("select item: %v", err) // %v: do not expose driver errors
	}
	return domain.Item{ID: row.ID, Name: row.Name}, nil
}
