-- +goose Up
-- +goose StatementBegin
INSERT INTO
    users (id, email, password)
VALUES
    (1, 'pedro@gmail.com', 'pedro123'),
    (2, 'fabio@gmail.com', 'fabio123'),
    (3, 'roger@gmail.com', 'roger123'),
    (4, 'joao@gmail.com', 'joao123');

INSERT INTO
    propostas (id, status, name, assignee_id)
VALUES
    (1, 'backlog', 'Eatopia - Garrafas e Guardanapos', 1),
    (2, 'backlog', 'Canelle - Itens variados', NULL),
    (3, 'backlog', 'Adega Santiago - Lacres', 1),
    (4, 'backlog', 'Dona Deola - Itens Variados', NULL);

INSERT INTO
    proposta_attachments (id, proposta_id, filename, mimetype)
VALUES
    (1, 1, 'orcamento_garrafas_2025.pdf', 'application/pdf'),
    (2, 1, 'guardanapos_personalizados.jpg', 'image/jpeg'),
    (3, 1, 'amostras_materiais.xlsx', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'),
    (4, 2, 'lista_produtos_canelle.pdf', 'application/pdf'),
    (5, 3, 'modelo_lacre_preview.png', 'image/png'),
    (6, 3, 'especificacoes_tecnicas.pdf', 'application/pdf'),
    (7, 3, 'orcamento_lacres_2025.xlsx', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE id IN (1, 2, 3, 4);
DELETE FROM propostas WHERE id IN (1, 2, 3, 4);
DELETE FROM proposta_attachments WHERE id IN (1, 2, 3, 4, 5, 6, 7);
-- +goose StatementEnd