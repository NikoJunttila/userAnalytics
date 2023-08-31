-- +goose Up 

CREATE TABLE visits (
  id TEXT NOT NULL,
  createdAt TIMESTAMP NOT NULL,
  country TEXT NOT NULL,
  ip TEXT NOT NULL,
  visitorStatus TEXT NOT NULL,
  domain TEXT NOT NULL,
  visitFrom TEXT NOT NULL
);

-- +goose Down
DROP TABLE visits;
