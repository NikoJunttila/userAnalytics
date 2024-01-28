-- name: CreateVisit :one
INSERT INTO visits(createdat,visitorstatus,visitDuration,domain,visitfrom)
VALUES($1,$2,$3,$4,$5)
RETURNING *;
--

-- name: GetTotalCount :one
SELECT
    COUNT(*) AS total_count,
    COUNT(CASE WHEN visitorstatus = 'new' THEN 1 END) AS new_visitor_count,
    CEIL(AVG(visitduration)) AS avg_visit_duration
FROM visits
WHERE domain = $1;
--

-- name: GetLimitedCount :many
SELECT
    COUNT(*) AS domain_count,
    COUNT(CASE WHEN visitorstatus = 'new' THEN 1 END) AS new_visitor_count,
    CEIL(AVG(visitduration)) AS avg_visit_duration,
    visitfrom, COUNT(*) AS count
FROM visits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY visitfrom;
--
-- name: GetSevenDays :many
SELECT
    COUNT(*) AS domain_count,
    COUNT(CASE WHEN visitorstatus = 'new' THEN 1 END) AS new_visitor_count,
    CEIL(AVG(visitduration)) AS avg_visit_duration,
    visitfrom, COUNT(*) AS count
FROM visits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '7 days'
GROUP BY visitfrom;
--
-- name: GetNinetyDays :many
SELECT
    COUNT(*) AS domain_count,
    COUNT(CASE WHEN visitorstatus = 'new' THEN 1 END) AS new_visitor_count,
    CEIL(AVG(visitduration)) AS avg_visit_duration,
    visitfrom, COUNT(*) AS count
FROM visits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '90 days'
GROUP BY visitfrom;
--
-- name: GetOsCount :many
SELECT
    COUNT(*) AS count,
    os AS column_value
FROM visits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL $2
GROUP BY os;
--
-- name: GetOsCount30 :many
SELECT
    COUNT(*) AS count,
    os AS column_value
FROM visits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY os;