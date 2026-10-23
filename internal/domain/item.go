// Package domain holds business entities and rules. It has no dependencies on
// storage, transport or frameworks.
package domain

import (
	"errors"
	"fmt"
	"strings"
)

// Domain errors. Callers compare them with errors.Is.
var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

// MaxNameLen is the maximum length of an item name in runes.
const MaxNameLen = 100

// Item is an example entity. Replace it with the service's own model.
type Item struct {
	ID   int64
	Name string
}

// NewItem validates input and returns an Item without an ID.
func NewItem(name string) (Item, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Item{}, fmt.Errorf("name is empty: %w", ErrInvalidInput)
	}
	if n := len([]rune(name)); n > MaxNameLen {
		return Item{}, fmt.Errorf("name has %d runes, max %d: %w", n, MaxNameLen, ErrInvalidInput)
	}
	return Item{Name: name}, nil
}
