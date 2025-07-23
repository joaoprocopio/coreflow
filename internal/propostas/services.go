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

type ListPropostasParams struct {
	Cursor int32 `json:"cursor"`
	Limit  int32 `json:"limit"`
}

func (s *Services) ListPropostas(ctx context.Context, cursor int32, limit int32) ([]*Proposta, error) {
	var propostas []*Proposta

	err := s.db.NewSelect().
		Model(&propostas).
		Relation("Assignee").
		Where("p.id > ?", cursor).
		Order("p.id ASC").
		Limit(int(limit)).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return propostas, nil
}

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

func (s *Services) ListPropostasWithAttachments(ctx context.Context, cursor int32, limit int32) ([]*Proposta, error) {
	var propostas []*Proposta

	err := s.db.NewSelect().
		Model(&propostas).
		Relation("Assignee").
		Relation("Attachments").
		Where("p.id > ?", cursor).
		Order("p.id ASC").
		Limit(int(limit)).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return propostas, nil
}
