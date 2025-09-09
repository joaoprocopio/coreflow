package services

import (
	"context"
	"coreflow/internal/db"

	"coreflow/gen/postgres/public/model"
	"coreflow/gen/postgres/public/table"

	"github.com/go-jet/jet/v2/postgres"
)

func New(db *db.DB) *Services {
	return &Services{db: db}
}

type Services struct {
	db *db.DB
}

// PropostaAssignee represents the assignee user data
type PropostaAssignee struct {
	ID    int32  `json:"id"`
	Email string `json:"email"`
}

// PropostaAttachment represents an attachment
type PropostaAttachment struct {
	ID       int32  `json:"id"`
	Mimetype string `json:"mimetype"`
	Filename string `json:"filename"`
}

// PropostaWithRelations represents a proposta with nested assignee and attachments
type PropostaWithRelations struct {
	ID          int32                `json:"id"`
	Status      model.PropostaStatus `json:"status"`
	Name        string               `json:"name"`
	Assignee    *PropostaAssignee    `json:"assignee"`
	Attachments []PropostaAttachment `json:"attachments"`
}

func (s *Services) ListPropostas(ctx context.Context, cursor int32, limit int32) ([]PropostaWithRelations, error) {
	// Create aliases for tables to avoid conflicts
	p := table.Propostas.AS("p")
	u := table.Users.AS("u")
	pa := table.PropostaAttachments.AS("pa")

	// We'll do a simpler approach: first get the paged propostas, then join with related data
	// This is equivalent to the CTE approach but simpler with Go Jet
	query := postgres.SELECT(
		p.ID.AS("proposta_id"),
		p.Status.AS("proposta_status"),
		p.Name.AS("proposta_name"),
		u.ID.AS("user_id"),
		u.Email.AS("user_email"),
		pa.ID.AS("attachment_id"),
		pa.Mimetype.AS("attachment_mimetype"),
		pa.Filename.AS("attachment_filename"),
	).FROM(
		p.LEFT_JOIN(u, u.ID.EQ(p.AssigneeID)).
			LEFT_JOIN(pa, pa.PropostaID.EQ(p.ID)),
	).WHERE(
		p.ID.GT(postgres.Int(int64(cursor))),
	).ORDER_BY(
		p.ID, pa.ID,
	) // Remove LIMIT here since we need to handle pagination differently

	// Generate SQL and args from Jet query
	sql, args := query.Sql()

	// Execute the query using pgx
	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process rows and group by proposta
	propostaMap := make(map[int32]*PropostaWithRelations)
	propostaOrder := []int32{} // To maintain order

	for rows.Next() {
		var (
			propostaID     int32
			propostaStatus model.PropostaStatus
			propostaName   string
			userID         *int32
			userEmail      *string
			attachmentID   *int32
			attachmentMime *string
			attachmentFile *string
		)

		err := rows.Scan(
			&propostaID, &propostaStatus, &propostaName,
			&userID, &userEmail,
			&attachmentID, &attachmentMime, &attachmentFile,
		)
		if err != nil {
			return nil, err
		}

		// Check if we've reached our limit of unique propostas
		if len(propostaMap) >= int(limit) && propostaMap[propostaID] == nil {
			break
		}

		// Get or create proposta
		proposta, exists := propostaMap[propostaID]
		if !exists {
			proposta = &PropostaWithRelations{
				ID:          propostaID,
				Status:      propostaStatus,
				Name:        propostaName,
				Attachments: []PropostaAttachment{},
			}

			// Set assignee if exists
			if userID != nil && userEmail != nil {
				proposta.Assignee = &PropostaAssignee{
					ID:    *userID,
					Email: *userEmail,
				}
			}

			propostaMap[propostaID] = proposta
			propostaOrder = append(propostaOrder, propostaID)
		}

		// Add attachment if exists and not already added
		if attachmentID != nil && attachmentMime != nil && attachmentFile != nil {
			// Check if attachment already exists (avoid duplicates)
			exists := false
			for _, existing := range proposta.Attachments {
				if existing.ID == *attachmentID {
					exists = true
					break
				}
			}

			if !exists {
				attachment := PropostaAttachment{
					ID:       *attachmentID,
					Mimetype: *attachmentMime,
					Filename: *attachmentFile,
				}
				proposta.Attachments = append(proposta.Attachments, attachment)
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Convert map to slice, maintaining order
	result := make([]PropostaWithRelations, 0, len(propostaOrder))
	for _, id := range propostaOrder {
		if proposta := propostaMap[id]; proposta != nil {
			result = append(result, *proposta)
		}
	}

	return result, nil
}
