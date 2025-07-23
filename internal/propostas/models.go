package propostas

import (
	"github.com/uptrace/bun"
)

type PropostaStatus string

const (
	PropostaStatusBacklog PropostaStatus = "backlog"
)

// User represents a user in the system
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID       int32  `bun:"id,pk,autoincrement" json:"id"`
	Email    string `bun:"email,notnull" json:"email"`
	Password string `bun:"password,notnull" json:"password"`
}

// Proposta represents a proposta/proposal
type Proposta struct {
	bun.BaseModel `bun:"table:propostas,alias:p"`

	ID         int32          `bun:"id,pk,autoincrement" json:"id"`
	Status     PropostaStatus `bun:"status,notnull" json:"status"`
	Name       string         `bun:"name,notnull" json:"name"`
	AssigneeID *int32         `bun:"assignee_id" json:"assignee_id"`

	// Relations
	Assignee    *User                 `bun:"rel:belongs-to,join:assignee_id=id" json:"assignee,omitempty"`
	Attachments []*PropostaAttachment `bun:"rel:has-many,join:id=proposta_id" json:"attachments,omitempty"`
}

// PropostaAttachment represents file attachments for propostas
type PropostaAttachment struct {
	bun.BaseModel `bun:"table:proposta_attachments,alias:pa"`

	ID         int32  `bun:"id,pk,autoincrement" json:"id"`
	PropostaID int32  `bun:"proposta_id,notnull" json:"proposta_id"`
	Filename   string `bun:"filename,notnull" json:"filename"`
	Mimetype   string `bun:"mimetype,notnull" json:"mimetype"`

	// Relations
	Proposta *Proposta `bun:"rel:belongs-to,join:proposta_id=id" json:"proposta,omitempty"`
}
