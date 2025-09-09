package services

import (
	"context"
	"coreflow/internal/db"

	"coreflow/gen/postgres/public/model"
	. "coreflow/gen/postgres/public/table"

	. "github.com/go-jet/jet/v2/postgres"
)

func New(db *db.DB) *Services {
	return &Services{db: db}
}

type Services struct {
	db *db.DB
}

func (s *Services) ListPropostas(ctx context.Context, cursor int32, limit int32) ([]model.Propostas, error) {
	query := SELECT(Propostas.ID, Propostas.Name, Propostas.Status, Propostas.AssigneeID).
		FROM(Propostas).
		WHERE(Propostas.ID.GT(Int(int64(cursor)))).
		ORDER_BY(Propostas.ID).
		LIMIT(int64(limit))

	// Generate SQL and args from Jet query
	sql, args := query.Sql()

	// Execute the query using pgx
	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var propostas []model.Propostas
	for rows.Next() {
		var p model.Propostas
		err := rows.Scan(&p.ID, &p.Name, &p.Status, &p.AssigneeID)
		if err != nil {
			return nil, err
		}
		propostas = append(propostas, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return propostas, nil
}
