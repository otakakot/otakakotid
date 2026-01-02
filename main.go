package main

import (
	"cmp"
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

	"github.com/otakakot/otakakotid/internal/handler"
	"github.com/otakakot/otakakotid/pkg/api"
)

func main() {
	port := cmp.Or(os.Getenv("PORT"), "8080")

	dsn := cmp.Or(
		os.Getenv("DSN"),
		"postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable",
	)

	conn, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), conn)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}

	hdl := api.HandlerFromMux(
		handler.New(pool),
		http.NewServeMux(),
	)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           hdl,
		ReadHeaderTimeout: 30 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	go func() {
		slog.Info("start server listen")

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-ctx.Done()

	slog.Info("start server shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}

	slog.Info("done server shutdown")
}
