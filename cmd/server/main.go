package main

import (
	"context"
	"convey/internal/config"
	"convey/internal/db"
	propostasQueries "convey/internal/propostas/queries"
	"convey/internal/server"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := run(ctx, logger); err != nil {
		logger.Error("error running server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	grp, ctx := errgroup.WithContext(ctx)

	cfg := config.New()

	db, err := db.New(ctx, cfg)

	if err != nil {
		return err
	}

	srv := server.NewServer(
		cfg,
		ctx,
		db,
		logger,
		propostasQueries.New(db),
	)

	grp.Go(func() error {
		logger.Info("server is listening", slog.String("address", fmt.Sprintf("http://%s", srv.Addr)))

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	grp.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		logger.Info("gracefully shutting down http server")

		if err := srv.Shutdown(ctx); err != nil {
			return err
		}

		return nil
	})

	grp.Go(func() error {
		<-ctx.Done()

		_, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		logger.Info("gracefully disconnecting from database")

		if err := db.Close(ctx); err != nil {
			return err
		}

		return nil
	})

	if err := grp.Wait(); err != nil {
		return err
	}

	return nil
}
