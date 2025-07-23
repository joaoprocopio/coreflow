package server

import (
	"context"
	"coreflow/internal/db"
	"coreflow/internal/health"
	"coreflow/internal/propostas"
	propostasQueries "coreflow/internal/propostas/queries"
	"log/slog"
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
	ctx context.Context,
	logger *slog.Logger,
	db *db.DB,
	propostasQueries *propostasQueries.Queries,
) {
	mux.Handle("GET /health", health.HandleHealth(ctx, logger, db))
	mux.Handle("GET /propostas", propostas.HandleListPropostas(ctx, logger, propostasQueries))
}
