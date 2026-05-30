package postgres

import (
	"context"
	"time"

	"neshiman/backend/db/sqlc"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CrossTeamRequestRepository struct {
	q    *sqlc.Queries
	pool *pgxpool.Pool
}

func NewCrossTeamRequestRepository(pool *pgxpool.Pool) ports.CrossTeamRequestRepository {
	return &CrossTeamRequestRepository{
		q:    sqlc.New(pool),
		pool: pool,
	}
}

func (r *CrossTeamRequestRepository) Create(ctx context.Context, req *domain.CrossTeamRequest) error {
	result, err := r.q.CreateCrossTeamRequest(ctx, sqlc.CreateCrossTeamRequestParams{
		RequestingUserID: req.RequestingUserID,
		TargetSeatID:     req.TargetSeatID,
		Date:             pgtype.Date{Time: time.Date(req.Date.Year, time.Month(req.Date.Month), req.Date.Day, 0, 0, 0, 0, time.UTC), Valid: true},
	})
	if err != nil {
		return err
	}
	req.ID = result.ID
	return nil
}

func (r *CrossTeamRequestRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.CrossTeamRequest, error) {
	result, err := r.q.GetCrossTeamRequestByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrCrossTeamRequestNotFound
		}
		return nil, err
	}
	return &domain.CrossTeamRequest{
		ID:               result.ID,
		RequestingUserID: result.RequestingUserID,
		TargetSeatID:     result.TargetSeatID,
		Date:             domain.Date{Year: result.Date.Time.Year(), Month: int(result.Date.Time.Month()), Day: result.Date.Time.Day()},
		Status:           domain.RequestStatus(result.Status),
	}, nil
}

func (r *CrossTeamRequestRepository) ListByStatus(ctx context.Context, status domain.RequestStatus) ([]domain.CrossTeamRequest, error) {
	results, err := r.q.ListCrossTeamRequestsByStatus(ctx, string(status))
	if err != nil {
		return nil, err
	}
	requests := make([]domain.CrossTeamRequest, len(results))
	for i, row := range results {
		requests[i] = domain.CrossTeamRequest{
			ID:               row.ID,
			RequestingUserID: row.RequestingUserID,
			TargetSeatID:     row.TargetSeatID,
			Date:             domain.Date{Year: row.Date.Time.Year(), Month: int(row.Date.Time.Month()), Day: row.Date.Time.Day()},
			Status:           domain.RequestStatus(row.Status),
		}
	}
	return requests, nil
}

func (r *CrossTeamRequestRepository) ListPendingByTeam(ctx context.Context, teamID uuid.UUID) ([]domain.CrossTeamRequest, error) {
	results, err := r.q.ListPendingRequestsForTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}
	requests := make([]domain.CrossTeamRequest, len(results))
	for i, row := range results {
		requests[i] = domain.CrossTeamRequest{
			ID:               row.ID,
			RequestingUserID: row.RequestingUserID,
			TargetSeatID:     row.TargetSeatID,
			Date:             domain.Date{Year: row.Date.Time.Year(), Month: int(row.Date.Time.Month()), Day: row.Date.Time.Day()},
			Status:           domain.RequestStatus(row.Status),
		}
	}
	return requests, nil
}

func (r *CrossTeamRequestRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.CrossTeamRequest, error) {
	results, err := r.q.ListCrossTeamRequestsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	requests := make([]domain.CrossTeamRequest, len(results))
	for i, row := range results {
		requests[i] = domain.CrossTeamRequest{
			ID:               row.ID,
			RequestingUserID: row.RequestingUserID,
			TargetSeatID:     row.TargetSeatID,
			Date:             domain.Date{Year: row.Date.Time.Year(), Month: int(row.Date.Time.Month()), Day: row.Date.Time.Day()},
			Status:           domain.RequestStatus(row.Status),
		}
	}
	return requests, nil
}

func (r *CrossTeamRequestRepository) ListAll(ctx context.Context) ([]domain.CrossTeamRequest, error) {
	results, err := r.q.ListCrossTeamRequestsAll(ctx)
	if err != nil {
		return nil, err
	}
	requests := make([]domain.CrossTeamRequest, len(results))
	for i, row := range results {
		requests[i] = domain.CrossTeamRequest{
			ID:               row.ID,
			RequestingUserID: row.RequestingUserID,
			TargetSeatID:     row.TargetSeatID,
			Date:             domain.Date{Year: row.Date.Time.Year(), Month: int(row.Date.Time.Month()), Day: row.Date.Time.Day()},
			Status:           domain.RequestStatus(row.Status),
		}
	}
	return requests, nil
}

func (r *CrossTeamRequestRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.RequestStatus) error {
	_, err := r.q.UpdateCrossTeamRequestStatus(ctx, sqlc.UpdateCrossTeamRequestStatusParams{
		ID:     id,
		Status: string(status),
	})
	return err
}
