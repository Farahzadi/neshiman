-- name: CreateSeat :one
INSERT INTO seats (room_id, team_id, label, pos_x, pos_y, assigned_user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetSeatByID :one
SELECT * FROM seats WHERE id = $1 AND deleted_at IS NULL;

-- name: ListSeatsByRoom :many
SELECT * FROM seats WHERE room_id = $1 AND deleted_at IS NULL ORDER BY pos_y, pos_x;

-- name: UpdateSeat :one
UPDATE seats
SET label = $2, pos_x = $3, pos_y = $4, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateSeatFull :one
UPDATE seats
SET label = $2, team_id = $3, pos_x = $4, pos_y = $5, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteSeat :exec
UPDATE seats SET deleted_at = now() WHERE id = $1;

-- name: SoftDeleteSeatsByRoom :exec
UPDATE seats SET deleted_at = now() WHERE room_id = $1;

-- name: GetSeatByAssignedUser :one
SELECT * FROM seats WHERE assigned_user_id = $1 AND deleted_at IS NULL;

-- name: AssignSeat :one
UPDATE seats SET assigned_user_id = $2, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UnassignSeat :one
UPDATE seats SET assigned_user_id = NULL, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;
