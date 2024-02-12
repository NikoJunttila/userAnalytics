-- name: CreatePageVisit :one
INSERT INTO pagevisits(createdat,domain,page)
VALUES($1,$2,$3)
RETURNING *;
--

-- name: GetPages7 :many
SELECT page, COUNT(*) as page_count
FROM pagevisits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '7 days'
GROUP BY page
ORDER BY page_count DESC;
--

-- name: GetPages30 :many
SELECT page, COUNT(*) as page_count
FROM pagevisits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY page
ORDER BY page_count DESC;
--
-- name: GetPages90 :many
SELECT page, COUNT(*) as page_count
FROM pagevisits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '90 days'
GROUP BY page
ORDER BY page_count DESC;
--
