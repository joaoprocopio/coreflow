package propostas

import (
	"context"
	"convey/internal/server/codec"
	"log/slog"
	"net/http"
	"strconv"
)

func HandleListPropostas(ctx context.Context, logger *slog.Logger, qrs *Queries) http.HandlerFunc {
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

		propostas, err := qrs.ListPropostas(ctx, ListPropostasParams{
			Cursor: cursor,
			Limit:  limit,
		})

		if err != nil {
			logger.Error("failed to list propostas", slog.String("error", err.Error()))
			http.Error(w, "failed to list propostas", http.StatusInternalServerError)
			return
		}

		ids := make([]int32, len(propostas))

		for i, proposta := range propostas {
			ids[i] = proposta.ID
		}

		attachments, err := qrs.ListPropostaAttachments(ctx, ids)

		if err != nil {
			logger.Error("failed to list proposta attachments", slog.String("error", err.Error()))
			http.Error(w, "failed to list proposta attachments", http.StatusInternalServerError)
			return
		}

		response := serializePropostasWithAttachments(propostas, attachments)

		codec.WriteEncodedJSON(w, r, http.StatusOK, response)
	}
}

type AssigneeResponse struct {
	ID    int32  `json:"id"`
	Email string `json:"email"`
}

type PropostaAttachmentResponse struct {
	ID       int32  `json:"id"`
	Filename string `json:"filename"`
	Mimetype string `json:"mimetype"`
}

type PropostaResponse struct {
	ID          int32                        `json:"id"`
	Status      string                       `json:"status"`
	Name        string                       `json:"name"`
	Assignee    *AssigneeResponse            `json:"assignee"`
	Attachments []PropostaAttachmentResponse `json:"attachments"`
}

func serializePropostasWithAttachments(propostas []ListPropostasRow, attachments []PropostaAttachment) []PropostaResponse {
	// Group attachments by proposta ID
	attachmentsByPropostaID := make(map[int32][]PropostaAttachmentResponse)
	for _, attachment := range attachments {
		attachmentsByPropostaID[attachment.PropostaID] = append(
			attachmentsByPropostaID[attachment.PropostaID],
			PropostaAttachmentResponse{
				ID:       attachment.ID,
				Filename: attachment.Filename,
				Mimetype: attachment.Mimetype,
			},
		)
	}

	response := make([]PropostaResponse, len(propostas))

	for i, proposta := range propostas {
		var assignee *AssigneeResponse

		// Check if assignee ID is valid (not null)
		if id, err := proposta.AssigneeID.Value(); err == nil && id != nil {
			assignee = &AssigneeResponse{
				ID:    int32(id.(int64)),
				Email: proposta.AssigneeEmail.String,
			}
		}

		// Get attachments for this proposta (empty slice if none)
		propostaAttachments := attachmentsByPropostaID[proposta.ID]
		if propostaAttachments == nil {
			propostaAttachments = []PropostaAttachmentResponse{}
		}

		response[i] = PropostaResponse{
			ID:          proposta.ID,
			Status:      string(proposta.Status),
			Name:        proposta.Name,
			Assignee:    assignee,
			Attachments: propostaAttachments,
		}
	}

	return response
}
