-- +goose Up 

CREATE TABLE visits (
  createdAt TIMESTAMP NOT NULL,
  visitorStatus TEXT NOT NULL,
  visitDuration INT NOT NULL,
  domain UUID NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
  visitFrom TEXT NOT NULL,
  browser TEXT NOT NULL,
  device TEXT NOT NULL,
  os TEXT NOT NULL,
  bounce BOOL
);

-- +goose Down
DROP TABLE visits;
