-- name: GetDomainsForUser :many
SELECT * FROM domain_follows WHERE user_id = ?;

-- name: CreateDomainFollow :exec
INSERT INTO domain_follows (id, created_at, user_id, domain_id, domain_name)
VALUES (?, ?, ?, ?, ?);
