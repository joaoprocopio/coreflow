-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(254) UNIQUE NOT NULL,
    password VARCHAR(128) NOT NULL
);

CREATE TYPE proposta_status AS ENUM('backlog');

CREATE TABLE propostas (
    id SERIAL PRIMARY KEY,
    status proposta_status NOT NULL,
    name VARCHAR(255) NOT NULL,
    assignee_id INTEGER NULL,

    FOREIGN KEY (assignee_id) REFERENCES users (id) ON DELETE SET NULL
);

CREATE TABLE proposta_attachments (
    id SERIAL PRIMARY KEY,
    proposta_id INTEGER NOT NULL,
    filename VARCHAR(255) NOT NULL,
    mimetype VARCHAR(255) NOT NULl,

    FOREIGN KEY (proposta_id) REFERENCES propostas (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE proposta_attachments;
DROP TABLE propostas;
DROP TYPE proposta_status;
-- +goose StatementEnd
