-- name: CreateReservation :one
INSERT INTO reservations (user_id, seat_id, date)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetReservationByID :one
SELECT * FROM reservations WHERE id = $1 AND deleted_at IS NULL;

-- name: GetReservationBySeatAndDate :one
SELECT * FROM reservations WHERE seat_id = $1 AND date = $2 AND deleted_at IS NULL;

-- name: ListReservationsByDate :many
SELECT * FROM reservations WHERE date = $1 AND deleted_at IS NULL ORDER BY created_at;

-- name: ListReservationsByUserAndDate :many
SELECT * FROM reservations WHERE user_id = $1 AND date = $2 AND deleted_at IS NULL ORDER BY created_at;

-- name: ListReservationsByUserAndWeek :many
SELECT * FROM reservations
WHERE user_id = $1
  AND date >= $2
  AND date <= $3
  AND deleted_at IS NULL
ORDER BY date;

-- name: CountReservationsByUserInWeek :one
SELECT COUNT(*)::int AS count
FROM reservations
WHERE user_id = $1
  AND date >= $2
  AND date <= $3
  AND deleted_at IS NULL;

-- name: ListReservationsByRoomAndDate :many
SELECT r.* FROM reservations r
JOIN seats s ON r.seat_id = s.id
WHERE s.room_id = $1 AND r.date = $2 AND r.deleted_at IS NULL
ORDER BY r.created_at;

-- name: ListReservationsByRoomAndDateWithDetails :many
SELECT r.id, r.user_id, u.name as user_name, r.seat_id, s.label as seat_label, r.date
FROM reservations r
JOIN seats s ON r.seat_id = s.id
JOIN users u ON r.user_id = u.id
WHERE s.room_id = $1 AND r.date = $2 AND r.deleted_at IS NULL
ORDER BY r.created_at;

-- name: SoftDeleteReservation :exec
UPDATE reservations SET deleted_at = now() WHERE id = $1;
