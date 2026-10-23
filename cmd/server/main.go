// Command server runs the HTTP API.
//
// Environment:
//
//	ADDR          listen address (default ":8080")
//	DATABASE_URL  PostgreSQL URL; empty means in-memory storage
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/7cout/go-template/internal/adapters/httpapi"
	"github.com/7cout/go-template/internal/adapters/memory"
	"github.com/7cout/go-template/internal/adapters/postgres"
	"github.com/7cout/go-template/internal/app"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	if err := run(log); err != nil {
		log.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	repo, closeRepo, err := newItemRepository(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer closeRepo()

	srv := &http.Server{
		Addr:              envOr("ADDR", ":8080"),
		Handler:           httpapi.NewHandler(app.NewItems(repo), log),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errc := make(chan error, 1)
	go func() {
		log.Info("listening", "addr", srv.Addr)
		errc <- srv.ListenAndServe()
	}()

	select {
	case err := <-errc:
		return fmt.Errorf("listen: %w", err)
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	if err := <-errc; !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen: %w", err)
	}
	return nil
}

func newItemRepository(ctx context.Context, dbURL string) (app.ItemRepository, func(), error) {
	if dbURL == "" {
		return memory.NewItems(), func() {}, nil
	}
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("connect postgres: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("ping postgres: %w", err)
	}
	return postgres.NewItems(pool), pool.Close, nil
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
