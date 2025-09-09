-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(254) UNIQUE NOT NULL,
    password VARCHAR(128) NOT NULL
);

CREATE TYPE task_status AS ENUM('backlog');

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    status task_status NOT NULL,
    name VARCHAR(255) NOT NULL,
    assignee_id INTEGER NULL,

    FOREIGN KEY (assignee_id) REFERENCES users (id) ON DELETE SET NULL
);

CREATE TABLE tasks_attachments (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL,
    filename VARCHAR(255) NOT NULL,
    mimetype VARCHAR(255) NOT NULl,

    FOREIGN KEY (task_id) REFERENCES tasks (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE tasks_attachments;
DROP TABLE tasks;
DROP TYPE task_status;
-- +goose StatementEnd
