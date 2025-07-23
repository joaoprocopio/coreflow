package health

import (
	"context"
	"convey/internal/db"
	"convey/internal/server/codec"
	"log/slog"
	"net/http"
)

func HandleHealth(ctx context.Context, logger *slog.Logger, db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type health struct {
			Server   string `json:"server"`
			Database string `json:"database"`
		}

		var err error

		err = db.Ping(ctx)

		if err != nil {
			logger.Error("database is not reachable", slog.String("error", err.Error()))
			http.Error(w, "database is not reachable", http.StatusInternalServerError)
			return
		}

		err = codec.WriteEncodedJSON(w, r, http.StatusOK, health{
			Server:   "ok",
			Database: "ok",
		})

		if err != nil {
			logger.Error("failed to write response", slog.String("error", err.Error()))
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	}

}
