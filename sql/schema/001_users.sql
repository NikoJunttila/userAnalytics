-- +goose Up 

CREATE TABLE users (
  id TEXT PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  passhash TEXT NOT NULL,
  api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT ''
);

-- +goose Down
DROP TABLE users;
