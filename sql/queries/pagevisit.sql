-- name: CreateVisit :one
INSERT INTO pagevisit(createdat,domain,page)
VALUES($1,$2,$3)
RETURNING *;
--

-- name: getPages :many
SELECT page, COUNT(*) as page_count
FROM pagevisits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL $2
GROUP BY page
ORDER BY page_count DESC;
--
-- name: getPages30 :many
SELECT page, COUNT(*) as page_count
FROM pagevisits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY page
ORDER BY page_count DESC;
--
-- name: getPages90 :many
SELECT page, COUNT(*) as page_count
FROM pagevisits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '90 days'
GROUP BY page
ORDER BY page_count DESC;
--
