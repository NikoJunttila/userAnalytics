-- +goose Up 

CREATE TABLE password_resets (
  id TEXT PRIMARY KEY,
  expiration TIMESTAMP NOT NULL,
  email TEXT NOT NULL,
  valid BOOL NOT NULL,
  token VARCHAR(128) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE password_resets;
