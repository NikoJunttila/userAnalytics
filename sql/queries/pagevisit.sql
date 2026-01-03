-- name: CreatePageVisit :exec
INSERT INTO pagevisits (createdAt, domain, page)
VALUES (?, ?, ?);

-- name: GetPages7 :many
SELECT page, COUNT(*) as page_count
FROM pagevisits
WHERE domain = ? AND createdAt >= date('now', '-7 days')
GROUP BY page
ORDER BY page_count DESC;

-- name: GetPages30 :many
SELECT page, COUNT(*) as page_count
FROM pagevisits
WHERE domain = ? AND createdAt >= date('now', '-30 days')
GROUP BY page
ORDER BY page_count DESC;

-- name: GetPages90 :many
SELECT page, COUNT(*) as page_count
FROM pagevisits
WHERE domain = ? AND createdAt >= date('now', '-90 days')
GROUP BY page
ORDER BY page_count DESC;
