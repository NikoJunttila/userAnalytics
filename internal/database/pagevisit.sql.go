// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: pagevisit.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPageVisit = `-- name: CreatePageVisit :one
INSERT INTO pagevisits(createdat,domain,page)
VALUES($1,$2,$3)
RETURNING createdat, domain, page
`

type CreatePageVisitParams struct {
	Createdat time.Time
	Domain    uuid.UUID
	Page      string
}

func (q *Queries) CreatePageVisit(ctx context.Context, arg CreatePageVisitParams) (Pagevisit, error) {
	row := q.db.QueryRowContext(ctx, createPageVisit, arg.Createdat, arg.Domain, arg.Page)
	var i Pagevisit
	err := row.Scan(&i.Createdat, &i.Domain, &i.Page)
	return i, err
}

const getPages = `-- name: GetPages :many

SELECT page, COUNT(*) as page_count
FROM pagevisits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL $2
GROUP BY page
ORDER BY page_count DESC
`

type GetPagesParams struct {
	Domain  uuid.UUID
	Column2 int64
}

type GetPagesRow struct {
	Page      string
	PageCount int64
}

func (q *Queries) GetPages(ctx context.Context, arg GetPagesParams) ([]GetPagesRow, error) {
	rows, err := q.db.QueryContext(ctx, getPages, arg.Domain, arg.Column2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPagesRow
	for rows.Next() {
		var i GetPagesRow
		if err := rows.Scan(&i.Page, &i.PageCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPages30 = `-- name: GetPages30 :many
SELECT page, COUNT(*) as page_count
FROM pagevisits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY page
ORDER BY page_count DESC
`

type GetPages30Row struct {
	Page      string
	PageCount int64
}

func (q *Queries) GetPages30(ctx context.Context, domain uuid.UUID) ([]GetPages30Row, error) {
	rows, err := q.db.QueryContext(ctx, getPages30, domain)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPages30Row
	for rows.Next() {
		var i GetPages30Row
		if err := rows.Scan(&i.Page, &i.PageCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPages90 = `-- name: GetPages90 :many
SELECT page, COUNT(*) as page_count
FROM pagevisits
WHERE domain = $1 AND createdat >= CURRENT_DATE - INTERVAL '90 days'
GROUP BY page
ORDER BY page_count DESC
`

type GetPages90Row struct {
	Page      string
	PageCount int64
}

func (q *Queries) GetPages90(ctx context.Context, domain uuid.UUID) ([]GetPages90Row, error) {
	rows, err := q.db.QueryContext(ctx, getPages90, domain)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPages90Row
	for rows.Next() {
		var i GetPages90Row
		if err := rows.Scan(&i.Page, &i.PageCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
