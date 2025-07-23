package propostas

import (
	"context"
	"convey/internal/server/codec"
	"log/slog"
	"net/http"
	"strconv"
)

type PropostaResponse struct {
	ID          int32                 `json:"id"`
	Status      PropostaStatus        `json:"status"`
	Name        string                `json:"name"`
	Assignee    *User                 `json:"assignee"`
	Attachments []*PropostaAttachment `json:"attachments"`
}

func HandleListPropostas(ctx context.Context, logger *slog.Logger, services *Services) http.HandlerFunc {
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

		// Use the simpler method that loads everything with relations
		propostas, err := services.ListPropostasWithAttachments(ctx, ListPropostasParams{
			Cursor: cursor,
			Limit:  limit,
		})

		if err != nil {
			logger.Error("failed to list propostas", slog.String("error", err.Error()))
			http.Error(w, "failed to list propostas", http.StatusInternalServerError)
			return
		}

		// Convert to response format
		response := make([]PropostaResponse, len(propostas))
		for i, proposta := range propostas {
			attachments := proposta.Attachments
			if attachments == nil {
				attachments = []*PropostaAttachment{}
			}

			response[i] = PropostaResponse{
				ID:          proposta.ID,
				Status:      proposta.Status,
				Name:        proposta.Name,
				Assignee:    proposta.Assignee,
				Attachments: attachments,
			}
		}

		codec.WriteEncodedJSON(w, r, http.StatusOK, response)
	}
}
