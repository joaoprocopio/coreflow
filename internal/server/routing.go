package server

import (
	"coreflow/internal/db"
	"coreflow/internal/health"
	"coreflow/internal/tasks"
	"log/slog"
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	db *db.DB,
	tasksSvc *tasks.Services,
) {
	mux.Handle("GET /health", health.HandleHealth(logger, db))
	mux.Handle("GET /api/v1/tasks", tasks.HandleListTasksV1(logger, tasksSvc))
}
