-- name: CreateVisit :one
INSERT INTO visits(createdat,visitorstatus,visitDuration,domain,visitfrom)
VALUES($1,$2,$3,$4,$5)
RETURNING *;
