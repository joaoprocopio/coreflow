package propostas

import (
	"context"
	"convey/internal/server/codec"
	"log/slog"
	"net/http"
	"strconv"
)

func HandleListPropostas(ctx context.Context, logger *slog.Logger, services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		var cursor int32 = 0
		var limit int32 = 10

		if c, err := strconv.ParseInt(params.Get("cursor"), 10, 32); err == nil {
			cursor = int32(c)
		}
		if l, err := strconv.ParseInt(params.Get("limit"), 10, 32); err == nil {
			limit = int32(l)
		}

		propostas, err := services.ListPropostasWithAttachments(ctx, cursor, limit)

		if err != nil {
			logger.Error("failed to list propostas", slog.String("error", err.Error()))
			http.Error(w, "failed to list propostas", http.StatusInternalServerError)
			return
		}

		serializedPropostas := SerializePropostasToResponse(propostas)

		codec.WriteEncodedJSON(w, r, http.StatusOK, serializedPropostas)
	}
}
