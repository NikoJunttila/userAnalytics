-- name: CreateVisit :one
INSERT INTO visits(createdat,visitorstatus,visitDuration,domain,visitfrom,browser,device,os,bounce)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)
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
-- name: GetOsCount7 :many
SELECT
    COUNT(*) AS count,
    os AS column_value
FROM visits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '7 days'
GROUP BY os;
--
-- name: GetOsCount30 :many
SELECT
    COUNT(*) AS count,
    os AS column_value
FROM visits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY os;
--
-- name: GetOsCount90 :many
SELECT
    COUNT(*) AS count,
    os AS column_value
FROM visits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '90 days'
GROUP BY os;
--
-- name: GetBrowserCount7 :many
SELECT
    COUNT(*) AS count,
    browser AS column_value
FROM visits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '7 days'
GROUP BY browser;
--

-- name: GetBrowserCount30 :many
SELECT
    COUNT(*) AS count,
    browser AS column_value
FROM visits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY browser;
--

-- name: GetBrowserCount90 :many
SELECT
    COUNT(*) AS count,
    browser AS column_value
FROM visits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '90 days'
GROUP BY browser;
--

-- name: GetDeviceCount7 :many
SELECT
    COUNT(*) AS count,
    device AS column_value
FROM visits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '7 days'
GROUP BY device;
--

-- name: GetDeviceCount30 :many
SELECT
    COUNT(*) AS count,
    device AS column_value
FROM visits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY device;
--

-- name: GetDeviceCount90 :many
SELECT
    COUNT(*) AS count,
    device AS column_value
FROM visits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '90 days'
GROUP BY device;
--
-- name: GetBounce7 :one
SELECT
	CEIL((COUNT(CASE WHEN bounce = true THEN 1 END) * 100.0 / COUNT(*))) AS bounced
FROM visits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '7 days';
--

-- name: GetBounce30 :one
SELECT
	CEIL((COUNT(CASE WHEN bounce = true THEN 1 END) * 100.0 / COUNT(*))) AS bounced
FROM visits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '30 days';
--

-- name: GetBounce90 :one
SELECT
	CEIL((COUNT(CASE WHEN bounce = true THEN 1 END) * 100.0 / COUNT(*))) AS bounced
FROM visits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '90 days';
--
