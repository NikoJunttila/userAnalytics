-- +goose Up
CREATE TABLE domain_follows (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    domain_id TEXT NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
    domain_name TEXT NOT NULL,
    UNIQUE (user_id, domain_id)
);

-- +goose Down
DROP TABLE domain_follows;
