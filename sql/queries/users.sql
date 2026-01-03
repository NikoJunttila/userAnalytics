-- name: CreateUser :exec
INSERT INTO users (id, created_at, updated_at, name, api_key, email, passhash)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE api_key = ?;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: UpdatePassword :exec
UPDATE users SET passhash = ? WHERE id = ?;

-- name: CreatePasswordReset :exec
INSERT INTO password_resets (id, expiration, email, valid, token)
VALUES (?, ?, ?, ?, ?);

-- name: ResetPassword :one
SELECT * FROM password_resets WHERE token = ?;

-- name: ResetInvalid :exec
UPDATE password_resets SET valid = 0 WHERE id = ?;
