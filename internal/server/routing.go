package server

import (
	"context"
	"coreflow/internal/db"
	"coreflow/internal/health"
	"coreflow/internal/propostas"
	propostasServices "coreflow/internal/propostas/services"
	"log/slog"
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
	ctx context.Context,
	logger *slog.Logger,
	db *db.DB,
	propostasServices *propostasServices.Services,
) {
	mux.Handle("GET /health", health.HandleHealth(ctx, logger, db))
	mux.Handle("GET /propostas", propostas.HandleListPropostas(ctx, logger, propostasServices))
}
