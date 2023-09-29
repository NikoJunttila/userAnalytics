-- +goose Up
CREATE TABLE domain_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    domain_id UUID NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
    domain_name TEXT NOT NULL,
    UNIQUE (user_id, domain_id)
);

-- +goose Down
DROP TABLE domain_follows;
