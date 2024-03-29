// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: domain_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createDomainFollow = `-- name: CreateDomainFollow :one

INSERT INTO domain_follows (id, created_at, user_id, domain_id, domain_name)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, user_id, domain_id, domain_name
`

type CreateDomainFollowParams struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UserID     uuid.UUID
	DomainID   uuid.UUID
	DomainName string
}

func (q *Queries) CreateDomainFollow(ctx context.Context, arg CreateDomainFollowParams) (DomainFollow, error) {
	row := q.db.QueryRowContext(ctx, createDomainFollow,
		arg.ID,
		arg.CreatedAt,
		arg.UserID,
		arg.DomainID,
		arg.DomainName,
	)
	var i DomainFollow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UserID,
		&i.DomainID,
		&i.DomainName,
	)
	return i, err
}

const getDomainsForUser = `-- name: GetDomainsForUser :many
SELECT id, created_at, user_id, domain_id, domain_name FROM domain_follows WHERE user_id = $1
`

func (q *Queries) GetDomainsForUser(ctx context.Context, userID uuid.UUID) ([]DomainFollow, error) {
	rows, err := q.db.QueryContext(ctx, getDomainsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DomainFollow
	for rows.Next() {
		var i DomainFollow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UserID,
			&i.DomainID,
			&i.DomainName,
		); err != nil {
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
