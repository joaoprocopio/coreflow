package tasks

import (
	"coreflow/internal/server/codec"
	"log/slog"
	"net/http"
	"strconv"
)

func HandleListTasksV1(logger *slog.Logger, svc *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		var cursor int32 = 0
		var limit int32 = 10

		if c, err := strconv.Atoi(params.Get("cursor")); err == nil {
			cursor = int32(c)
		}
		if l, err := strconv.Atoi(params.Get("limit")); err == nil {
			limit = int32(l)
		}

		tasks, err := svc.ListTasks(r.Context(), cursor, limit)

		if err != nil {
			logger.Error("failed to list tasks", slog.String("error", err.Error()))
			http.Error(w, "failed to list tasks", http.StatusInternalServerError)
			return
		}

		codec.WriteEncodedJSON(w, r, http.StatusOK, tasks)
	}
}
