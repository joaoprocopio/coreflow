package propostas

import (
	"context"
	"convey/internal/database"

	"github.com/uptrace/bun"
)

type Services struct {
	db *database.DB
}

func NewServices(db *database.DB) *Services {
	return &Services{db: db}
}

// ListPropostasParams contains parameters for listing propostas
type ListPropostasParams struct {
	Cursor int32 `json:"cursor"`
	Limit  int32 `json:"limit"`
}

// ListPropostas retrieves a paginated list of propostas with assignee information
func (s *Services) ListPropostas(ctx context.Context, params ListPropostasParams) ([]*Proposta, error) {
	var propostas []*Proposta

	err := s.db.NewSelect().
		Model(&propostas).
		Relation("Assignee").
		Where("p.id > ?", params.Cursor).
		Order("p.id ASC").
		Limit(int(params.Limit)).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return propostas, nil
}

// ListPropostaAttachments retrieves all attachments for given proposta IDs
func (s *Services) ListPropostaAttachments(ctx context.Context, propostaIDs []int32) ([]*PropostaAttachment, error) {
	var attachments []*PropostaAttachment

	err := s.db.NewSelect().
		Model(&attachments).
		Where("pa.proposta_id IN (?)", bun.In(propostaIDs)).
		Order("pa.proposta_id ASC", "pa.id ASC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return attachments, nil
}

// ListPropostasWithAttachments retrieves propostas with their attachments in a single query
func (s *Services) ListPropostasWithAttachments(ctx context.Context, params ListPropostasParams) ([]*Proposta, error) {
	var propostas []*Proposta

	err := s.db.NewSelect().
		Model(&propostas).
		Relation("Assignee").
		Relation("Attachments").
		Where("p.id > ?", params.Cursor).
		Order("p.id ASC").
		Limit(int(params.Limit)).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return propostas, nil
}
