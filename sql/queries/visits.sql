-- name: CreateVisit :one
INSERT INTO visits(id,createdat, country,ip,visitorstatus,domain,visitfrom)
VALUES($1,$2,$3,$4,$5,$6,$7)
RETURNING *;
