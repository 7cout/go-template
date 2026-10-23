package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/7cout/go-template/internal/adapters/memory"
	"github.com/7cout/go-template/internal/app"
	"github.com/7cout/go-template/internal/domain"
)

func TestItems_CreateAndGet(t *testing.T) {
	t.Parallel()
	svc := app.NewItems(memory.NewItems())

	created, err := svc.Create(t.Context(), "apple")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	got, err := svc.Get(t.Context(), created.ID)
	if err != nil {
		t.Fatalf("Get(%d) error = %v", created.ID, err)
	}
	if got != created {
		t.Errorf("Get(%d) = %+v, want %+v", created.ID, got, created)
	}
}

func TestItems_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		run     func(ctx context.Context, svc *app.Items) error
		wantErr error
	}{
		{
			name: "create invalid",
			run: func(ctx context.Context, svc *app.Items) error {
				_, err := svc.Create(ctx, "")
				return err
			},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "get missing",
			run: func(ctx context.Context, svc *app.Items) error {
				_, err := svc.Get(ctx, 42)
				return err
			},
			wantErr: domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.run(t.Context(), app.NewItems(memory.NewItems()))
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
