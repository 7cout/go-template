// Package httpapi exposes use cases over HTTP/JSON.
package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/7cout/go-template/internal/domain"
)

// ItemService is what the handler needs from the app layer.
type ItemService interface {
	Create(ctx context.Context, name string) (domain.Item, error)
	Get(ctx context.Context, id int64) (domain.Item, error)
}

// Handler routes HTTP requests to the item service.
type Handler struct {
	items ItemService
	log   *slog.Logger
	mux   *http.ServeMux
}

// NewHandler registers routes and returns a ready http.Handler.
func NewHandler(items ItemService, log *slog.Logger) *Handler {
	h := &Handler{items: items, log: log, mux: http.NewServeMux()}
	h.mux.HandleFunc("GET /healthz", h.health)
	h.mux.HandleFunc("POST /items", h.createItem)
	h.mux.HandleFunc("GET /items/{id}", h.getItem)
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

type itemDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createItem(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, "malformed JSON")
		return
	}
	item, err := h.items.Create(r.Context(), req.Name)
	if err != nil {
		h.writeDomainError(w, r, err)
		return
	}
	h.writeJSON(w, r, http.StatusCreated, itemDTO{ID: item.ID, Name: item.Name})
}

func (h *Handler) getItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.writeError(w, r, http.StatusBadRequest, "id must be an integer")
		return
	}
	item, err := h.items.Get(r.Context(), id)
	if err != nil {
		h.writeDomainError(w, r, err)
		return
	}
	h.writeJSON(w, r, http.StatusOK, itemDTO{ID: item.ID, Name: item.Name})
}

func (h *Handler) writeDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		h.writeError(w, r, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		h.writeError(w, r, http.StatusNotFound, "not found")
	default:
		h.log.ErrorContext(r.Context(), "request failed", "path", r.URL.Path, "err", err)
		h.writeError(w, r, http.StatusInternalServerError, "internal error")
	}
}

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, status int, msg string) {
	h.writeJSON(w, r, status, map[string]string{"error": msg})
}

func (h *Handler) writeJSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.log.WarnContext(r.Context(), "write response", "err", err)
	}
}
