-- name: CreateUser :one
INSERT INTO users (name, email, team_id, role, weekly_limit)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ListUsersByTeam :many
SELECT * FROM users WHERE team_id = $1 ORDER BY name;

-- name: ListUsers :many
SELECT * FROM users ORDER BY name;

-- name: UpdateUserWeeklyLimit :one
UPDATE users SET weekly_limit = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
