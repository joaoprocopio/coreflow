-- name: ListPropostas :many
WITH paged_propostas AS (
    SELECT *
    FROM propostas AS p
    WHERE
        p.id > sqlc.arg('cursor')
    ORDER BY p.id
    LIMIT sqlc.arg('limit')
)

SELECT
    p.id,
    p.status,
    p.name,
    CASE
        WHEN u.id IS NULL THEN NULL
        ELSE json_build_object(
            'id', u.id,
            'email', u.email
        )
    END AS assignee,
    COALESCE(json_agg(
        json_build_object(
            'id', pa.id,
            'mimetype', pa.mimetype,
            'filename', pa.filename
        )
    ) FILTER (WHERE pa.id IS NOT NULL), '[]') AS attachments
FROM paged_propostas AS p

LEFT JOIN users AS u
    ON u.id = p.assignee_id
LEFT JOIN proposta_attachments AS pa
    ON pa.proposta_id = p.id

GROUP BY p.id, p.status, p.name, u.id, u.email;
