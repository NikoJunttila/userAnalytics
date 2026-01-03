-- name: CreateDomain :exec
INSERT INTO domains (id, owner_id, name, url, total_visits, total_unique, total_time, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetDomains :many
SELECT * FROM domains;

-- name: GetDomain :one
SELECT * FROM domains WHERE id = ?;

-- name: UpdateDomain :exec
UPDATE domains SET 
  total_visits = total_visits + ?,
  total_unique = total_unique + ?
WHERE id = ?;

-- name: GetPrevMonthStats :one
SELECT
  COUNT(*) AS total_count,
  COUNT(CASE WHEN visitorStatus = 'new' THEN 1 END) AS new_visitor_count
FROM visits 
WHERE domain = ?
  AND createdAt >= date('now', '-60 days')
  AND createdAt < date('now', '-30 days');

-- name: GetMonthStats :one
SELECT
  COUNT(*) AS total_count,
  COUNT(CASE WHEN visitorStatus = 'new' THEN 1 END) AS new_visitor_count
FROM visits
WHERE domain = ?
  AND createdAt >= date('now', '-30 days');
