-- name: CreateSeat :one
INSERT INTO seats (room_id, team_id, label, pos_x, pos_y, rotation)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetSeatByID :one
SELECT * FROM seats WHERE id = $1;

-- name: ListSeatsByRoom :many
SELECT * FROM seats WHERE room_id = $1 ORDER BY pos_y, pos_x;

-- name: UpdateSeat :one
UPDATE seats
SET label = $2, pos_x = $3, pos_y = $4, rotation = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateSeatFull :one
UPDATE seats
SET label = $2, team_id = $3, pos_x = $4, pos_y = $5, rotation = $6, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteSeat :exec
DELETE FROM seats WHERE id = $1;
