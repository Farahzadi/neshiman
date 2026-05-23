-- name: CreateUser :one
INSERT INTO users (name, email, team_id, role, weekly_limit, password_hash)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserByName :one
SELECT * FROM users WHERE name = $1 AND deleted_at IS NULL;

-- name: ListUsersByTeam :many
SELECT * FROM users WHERE team_id = $1 AND deleted_at IS NULL ORDER BY name;

-- name: ListUsers :many
SELECT * FROM users WHERE deleted_at IS NULL ORDER BY name;

-- name: UpdateUserWeeklyLimit :one
UPDATE users SET weekly_limit = $2, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users SET password_hash = $2, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteUser :exec
UPDATE users SET deleted_at = now() WHERE id = $1;
