-- name: GetDomainsForUser :many
SELECT * FROM domain_follows WHERE user_id = $1;
--

-- name: CreateDomainFollow :one
INSERT INTO domain_follows (id, created_at, user_id, domain_id)
VALUES ($1, $2, $3, $4)
RETURNING *;
--
