-- name: ListPropostas :many
SELECT 
    p.id,
    p.status,
    p.name,
    u.id AS assignee_id,
    u.email AS assignee_email
FROM propostas AS p
LEFT JOIN users AS u
    ON u.id = p.assignee_id
WHERE p.id > sqlc.arg('cursor')
ORDER BY p.id ASC
LIMIT sqlc.arg('limit');

-- name: ListPropostaAttachments :many
SELECT 
    pa.id,
    pa.proposta_id,
    pa.filename,
    pa.mimetype
FROM proposta_attachments AS pa
WHERE pa.proposta_id = ANY(sqlc.arg('proposta_ids')::int[]);
