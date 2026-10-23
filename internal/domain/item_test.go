package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/7cout/go-template/internal/domain"
)

func TestNewItem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		in       string
		wantName string
		wantErr  error
	}{
		{name: "ok", in: "apple", wantName: "apple"},
		{name: "trims spaces", in: "  apple  ", wantName: "apple"},
		{name: "cyrillic within limit", in: strings.Repeat("я", domain.MaxNameLen), wantName: strings.Repeat("я", domain.MaxNameLen)},
		{name: "empty", in: "", wantErr: domain.ErrInvalidInput},
		{name: "only spaces", in: "   ", wantErr: domain.ErrInvalidInput},
		{name: "too long", in: strings.Repeat("a", domain.MaxNameLen+1), wantErr: domain.ErrInvalidInput},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := domain.NewItem(tt.in)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewItem(%q) error = %v, want %v", tt.in, err, tt.wantErr)
			}
			if got.Name != tt.wantName {
				t.Errorf("NewItem(%q).Name = %q, want %q", tt.in, got.Name, tt.wantName)
			}
		})
	}
}
