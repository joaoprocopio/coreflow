package server

import (
	"context"
	"coreflow/internal/config"
	"coreflow/internal/db"
	"coreflow/internal/propostas"
	"log/slog"
	"net"
	"net/http"
)

func NewServer(
	cfg *config.Config,
	ctx context.Context,
	db *db.DB,
	logger *slog.Logger,
	propostasSvc *propostas.Services,
) *http.Server {
	var mux *http.ServeMux = http.NewServeMux()

	addRoutes(
		mux,
		logger,
		db,
		propostasSvc,
	)

	var handler http.Handler = mux

	handler = loggerMiddleware(handler, logger)

	var srv *http.Server = &http.Server{
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
		Addr:     net.JoinHostPort(cfg.SrvHost, cfg.SrvPort),
		Handler:  handler,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	return srv
}
