-- name: CreateVisit :exec
INSERT INTO visits (createdAt, visitorStatus, visitDuration, domain, visitFrom, browser, device, os, bounce)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetTotalCount :one
SELECT
    COUNT(*) AS total_count,
    COUNT(CASE WHEN visitorStatus = 'new' THEN 1 END) AS new_visitor_count,
    AVG(visitDuration) AS avg_visit_duration
FROM visits
WHERE domain = ?;

-- name: GetLimitedCount :many
SELECT
    COUNT(*) AS domain_count,
    COUNT(CASE WHEN visitorStatus = 'new' THEN 1 END) AS new_visitor_count,
    AVG(visitDuration) AS avg_visit_duration,
    visitFrom,
    COUNT(*) AS count
FROM visits
WHERE domain = ? AND createdAt >= date('now', '-30 days')
GROUP BY visitFrom;

-- name: GetSevenDays :many
SELECT
    COUNT(*) AS domain_count,
    COUNT(CASE WHEN visitorStatus = 'new' THEN 1 END) AS new_visitor_count,
    AVG(visitDuration) AS avg_visit_duration,
    visitFrom,
    COUNT(*) AS count
FROM visits
WHERE domain = ? AND createdAt >= date('now', '-7 days')
GROUP BY visitFrom;

-- name: GetNinetyDays :many
SELECT
    COUNT(*) AS domain_count,
    COUNT(CASE WHEN visitorStatus = 'new' THEN 1 END) AS new_visitor_count,
    AVG(visitDuration) AS avg_visit_duration,
    visitFrom,
    COUNT(*) AS count
FROM visits
WHERE domain = ? AND createdAt >= date('now', '-90 days')
GROUP BY visitFrom;

-- name: GetOsCount7 :many
SELECT
    COUNT(*) AS count,
    os AS column_value
FROM visits
WHERE domain = ? AND createdAt >= date('now', '-7 days')
GROUP BY os;

-- name: GetOsCount30 :many
SELECT
    COUNT(*) AS count,
    os AS column_value
FROM visits
WHERE domain = ? AND createdAt >= date('now', '-30 days')
GROUP BY os;

-- name: GetOsCount90 :many
SELECT
    COUNT(*) AS count,
    os AS column_value
FROM visits
WHERE domain = ? AND createdAt >= date('now', '-90 days')
GROUP BY os;

-- name: GetBrowserCount7 :many
SELECT
    COUNT(*) AS count,
    browser AS column_value
FROM visits
WHERE domain = ? AND createdAt >= date('now', '-7 days')
GROUP BY browser;

-- name: GetBrowserCount30 :many
SELECT
    COUNT(*) AS count,
    browser AS column_value
FROM visits
WHERE domain = ? AND createdAt >= date('now', '-30 days')
GROUP BY browser;

-- name: GetBrowserCount90 :many
SELECT
    COUNT(*) AS count,
    browser AS column_value
FROM visits
WHERE domain = ? AND createdAt >= date('now', '-90 days')
GROUP BY browser;

-- name: GetDeviceCount7 :many
SELECT
    COUNT(*) AS count,
    device AS column_value
FROM visits
WHERE domain = ? AND createdAt >= date('now', '-7 days')
GROUP BY device;

-- name: GetDeviceCount30 :many
SELECT
    COUNT(*) AS count,
    device AS column_value
FROM visits
WHERE domain = ? AND createdAt >= date('now', '-30 days')
GROUP BY device;

-- name: GetDeviceCount90 :many
SELECT
    COUNT(*) AS count,
    device AS column_value
FROM visits
WHERE domain = ? AND createdAt >= date('now', '-90 days')
GROUP BY device;

-- name: GetBounce7 :one
SELECT
    (COUNT(CASE WHEN bounce = 1 THEN 1 END) * 100.0 / COUNT(*)) AS bounced
FROM visits
WHERE domain = ? AND createdAt >= date('now', '-7 days');

-- name: GetBounce30 :one
SELECT
    (COUNT(CASE WHEN bounce = 1 THEN 1 END) * 100.0 / COUNT(*)) AS bounced
FROM visits
WHERE domain = ? AND createdAt >= date('now', '-30 days');

-- name: GetBounce90 :one
SELECT
    (COUNT(CASE WHEN bounce = 1 THEN 1 END) * 100.0 / COUNT(*)) AS bounced
FROM visits
WHERE domain = ? AND createdAt >= date('now', '-90 days');
