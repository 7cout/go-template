package memory_test

import (
	"sync"
	"testing"

	"github.com/7cout/go-template/internal/adapters/memory"
	"github.com/7cout/go-template/internal/domain"
)

func TestItems_ConcurrentCreate(t *testing.T) {
	t.Parallel()
	repo := memory.NewItems()

	const n = 50
	var wg sync.WaitGroup
	for range n {
		wg.Go(func() {
			if _, err := repo.Create(t.Context(), domain.Item{Name: "x"}); err != nil {
				t.Errorf("Create() error = %v", err)
			}
		})
	}
	wg.Wait()

	if _, err := repo.Get(t.Context(), n); err != nil {
		t.Errorf("Get(%d) error = %v, want item after %d creates", n, err, n)
	}
}
