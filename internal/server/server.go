package server

import (
	"context"
	"coreflow/internal/config"
	"coreflow/internal/db"
	propostasServices "coreflow/internal/propostas/services"
	"log/slog"
	"net"
	"net/http"
)

func NewServer(
	cfg *config.Config,
	ctx context.Context,
	db *db.DB,
	logger *slog.Logger,
	propostasServices *propostasServices.Services,
) *http.Server {
	var mux *http.ServeMux = http.NewServeMux()

	addRoutes(
		mux,
		ctx,
		logger,
		db,
		propostasServices,
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
