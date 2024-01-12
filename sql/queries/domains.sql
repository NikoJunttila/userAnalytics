-- name: CreateDomain :one
INSERT INTO domains (id, owner_id,name,url,total_visits,total_unique,total_time,created_at,updated_at)
VALUES ($1, $2, $3, $4, $5, $6,$7,$8,$9)
RETURNING *;

-- name: GetDomains :many
SELECT * FROM domains;
--

-- name: GetDomain :one
Select * FROM domains WHERE id = $1;
--

-- name: UpdateDomain :exec
UPDATE domains
  set total_visits = total_visits + $2,
  total_unique =  total_unique + $3
WHERE id = $1;
--

-- name: GetPrevMonthStats :one
SELECT
  COUNT(*) AS total_count,
  COUNT(CASE WHEN visitorstatus='new' THEN 1 END) AS new_visitor_count
FROM visits 
WHERE domain = $1
  AND createdat >= CURRENT_DATE - INTERVAL '60 days'
  AND createdat < CURRENT_DATE - INTERVAL '30 days';
--

-- name: GetMonthStats :one
SELECT
  COUNT(*) AS total_count,
  COUNT(CASE WHEN visitorstatus='new' THEN 1 END) AS new_visitor_count
FROM visits
WHERE domain = $1
  AND createdat >= CURRENT_DATE - INTERVAL '30 days';
