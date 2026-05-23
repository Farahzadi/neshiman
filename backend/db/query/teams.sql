-- name: CreateTeam :one
INSERT INTO teams (name)
VALUES ($1)
RETURNING *;

-- name: GetTeamByID :one
SELECT * FROM teams WHERE id = $1 AND deleted_at IS NULL;

-- name: ListTeams :many
SELECT * FROM teams WHERE deleted_at IS NULL ORDER BY name;

-- name: SoftDeleteTeam :exec
UPDATE teams SET deleted_at = now() WHERE id = $1;
