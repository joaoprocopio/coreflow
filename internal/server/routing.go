package server

import (
	"coreflow/internal/db"
	"coreflow/internal/health"
	"coreflow/internal/propostas"
	"log/slog"
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	db *db.DB,
	propostasSvc *propostas.Services,
) {
	mux.Handle("GET /health", health.HandleHealth(logger, db))
	mux.Handle("GET /propostas", propostas.HandleListPropostas(logger, propostasSvc))
}
