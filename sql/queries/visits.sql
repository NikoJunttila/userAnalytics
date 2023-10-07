-- name: CreateVisit :one
INSERT INTO visits(createdat,visitorstatus,visitDuration,domain,visitfrom)
VALUES($1,$2,$3,$4,$5)
RETURNING *;
--

-- name: GetTotalCount :one
SELECT
    (SELECT COUNT(*) FROM visits) AS total_count,
    (SELECT COUNT(visitorstatus) FROM visits WHERE visitorstatus = 'new') AS new_visitor_count,
    (SELECT CEIL(AVG(visitduration)) FROM visits ) AS avg_visit_duration
FROM visits WHERE visits.domain = $1
GROUP BY total_count;
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
