-- name: CreateRoom :one
INSERT INTO rooms (name, grid_width, grid_height)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRoomByID :one
SELECT * FROM rooms WHERE id = $1;

-- name: ListRooms :many
SELECT * FROM rooms ORDER BY created_at DESC;

-- name: UpdateRoom :one
UPDATE rooms
SET name = $2, grid_width = $3, grid_height = $4, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteRoom :exec
DELETE FROM rooms WHERE id = $1;
