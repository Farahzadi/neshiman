-- name: CreateCrossTeamRequest :one
INSERT INTO cross_team_requests (requesting_user_id, target_seat_id, date)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCrossTeamRequestByID :one
SELECT * FROM cross_team_requests WHERE id = $1 AND deleted_at IS NULL;

-- name: ListCrossTeamRequestsByStatus :many
SELECT * FROM cross_team_requests WHERE status = $1 AND deleted_at IS NULL ORDER BY created_at DESC;

-- name: ListPendingRequestsForTeam :many
SELECT ctr.* FROM cross_team_requests ctr
JOIN seats s ON ctr.target_seat_id = s.id
JOIN teams t ON s.team_id = t.id
WHERE t.id = $1 AND ctr.status = 'pending' AND ctr.deleted_at IS NULL
ORDER BY ctr.created_at DESC;

-- name: UpdateCrossTeamRequestStatus :one
UPDATE cross_team_requests
SET status = $2, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;
