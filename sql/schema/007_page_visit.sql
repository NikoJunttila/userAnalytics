-- +goose Up 

CREATE TABLE pagevisits (
  createdAt TIMESTAMP NOT NULL,
  domain UUID NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
  page TEXT NOT NULL
);

-- +goose Down
DROP TABLE pagevisits;
