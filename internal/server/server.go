package server

import (
	"context"
	"convey/internal/config"
	"convey/internal/db"
	propostasQueries "convey/internal/propostas/queries"
	"log/slog"
	"net"
	"net/http"
)

func NewServer(cfg *config.Config, ctx context.Context, db *db.DB, logger *slog.Logger, propostasQueries *propostasQueries.Queries) *http.Server {
	var mux *http.ServeMux = http.NewServeMux()

	addRoutes(
		mux,
		ctx,
		logger,
		db,
		propostasQueries,
	)

	var handler http.Handler = mux

	handler = loggerMiddleware(handler, logger)

	var srv *http.Server = &http.Server{
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
		Addr:     net.JoinHostPort(cfg.SrvHost, cfg.SrvPort),
		Handler:  handler,
	}

	return srv
}
