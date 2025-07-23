package server

import (
	"context"
	"convey/internal/database"
	"convey/internal/health"
	"convey/internal/propostas"
	"log/slog"
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
	ctx context.Context,
	logger *slog.Logger,
	db *database.DB,
	propostasServices *propostas.Services,
) {
	mux.Handle("GET /health", health.HandleHealth(ctx, logger, db))
	mux.Handle("GET /propostas", propostas.HandleListPropostas(ctx, logger, propostasServices))
}
