-- name: CreateReservation :one
INSERT INTO reservations (user_id, seat_id, date)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetReservationByID :one
SELECT * FROM reservations WHERE id = $1;

-- name: GetReservationBySeatAndDate :one
SELECT * FROM reservations WHERE seat_id = $1 AND date = $2;

-- name: ListReservationsByUserAndWeek :many
SELECT * FROM reservations
WHERE user_id = $1
  AND date >= $2
  AND date <= $3
ORDER BY date;

-- name: CountReservationsByUserInWeek :one
SELECT COUNT(*)::int AS count
FROM reservations
WHERE user_id = $1
  AND date >= $2
  AND date <= $3;

-- name: DeleteReservation :exec
DELETE FROM reservations WHERE id = $1;
