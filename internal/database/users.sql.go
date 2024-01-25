// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPasswordReset = `-- name: CreatePasswordReset :one
INSERT INTO password_resets(id,expiration,email,valid,token)
VALUES($1,$2,$3,$4,$5)
RETURNING id, expiration, email, valid, token
`

type CreatePasswordResetParams struct {
	ID         uuid.UUID
	Expiration time.Time
	Email      string
	Valid      bool
	Token      string
}

func (q *Queries) CreatePasswordReset(ctx context.Context, arg CreatePasswordResetParams) (PasswordReset, error) {
	row := q.db.QueryRowContext(ctx, createPasswordReset,
		arg.ID,
		arg.Expiration,
		arg.Email,
		arg.Valid,
		arg.Token,
	)
	var i PasswordReset
	err := row.Scan(
		&i.ID,
		&i.Expiration,
		&i.Email,
		&i.Valid,
		&i.Token,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users(id,created_at,updated_at,name, api_key,email,passhash)
VALUES($1,$2,$3,$4,
encode(sha256(random()::text::bytea), 'hex'),$5,$6
)
RETURNING id, created_at, updated_at, name, email, passhash, api_key
`

type CreateUserParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Email     string
	Passhash  string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Email,
		arg.Passhash,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Email,
		&i.Passhash,
		&i.ApiKey,
	)
	return i, err
}

const getUserByAPIKey = `-- name: GetUserByAPIKey :one
SELECT id, created_at, updated_at, name, email, passhash, api_key FROM users WHERE api_key = $1
`

func (q *Queries) GetUserByAPIKey(ctx context.Context, apiKey string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByAPIKey, apiKey)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Email,
		&i.Passhash,
		&i.ApiKey,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one

SELECT id, created_at, updated_at, name, email, passhash, api_key FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Email,
		&i.Passhash,
		&i.ApiKey,
	)
	return i, err
}

const resetPassword = `-- name: ResetPassword :one
SELECT id, expiration, email, valid, token FROM password_resets WHERE token = $1
`

func (q *Queries) ResetPassword(ctx context.Context, token string) (PasswordReset, error) {
	row := q.db.QueryRowContext(ctx, resetPassword, token)
	var i PasswordReset
	err := row.Scan(
		&i.ID,
		&i.Expiration,
		&i.Email,
		&i.Valid,
		&i.Token,
	)
	return i, err
}

const updatePassword = `-- name: UpdatePassword :exec
UPDATE users
SET passhash = $1
WHERE id = $2
`

type UpdatePasswordParams struct {
	Passhash string
	ID       uuid.UUID
}

func (q *Queries) UpdatePassword(ctx context.Context, arg UpdatePasswordParams) error {
	_, err := q.db.ExecContext(ctx, updatePassword, arg.Passhash, arg.ID)
	return err
}
