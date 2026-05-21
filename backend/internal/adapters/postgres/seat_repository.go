package postgres

import (
	"context"
	"neshiman/backend/db/sqlc"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeatRepository struct {
	q    *sqlc.Queries
	pool *pgxpool.Pool
}

func NewSeatRepository(pool *pgxpool.Pool) ports.SeatRepository {
	return &SeatRepository{
		q:    sqlc.New(pool),
		pool: pool,
	}
}

func (r *SeatRepository) Create(ctx context.Context, seat *domain.Seat) error {
	result, err := r.q.CreateSeat(ctx, sqlc.CreateSeatParams{
		RoomID:   seat.RoomID,
		TeamID:   seat.TeamID,
		Label:    seat.Label,
		PosX:     int32(seat.Position.X),
		PosY:     int32(seat.Position.Y),
		Rotation: int32(seat.Rotation),
	})
	if err != nil {
		return err
	}
	seat.ID = result.ID
	return nil
}

func (r *SeatRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Seat, error) {
	result, err := r.q.GetSeatByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrSeatNotFound
		}
		return nil, err
	}
	return &domain.Seat{
		ID:     result.ID,
		RoomID: result.RoomID,
		TeamID: result.TeamID,
		Label:  result.Label,
		Position: domain.Position{
			X: int(result.PosX),
			Y: int(result.PosY),
		},
		Rotation: domain.Rotation(result.Rotation),
	}, nil
}

func (r *SeatRepository) ListByRoom(ctx context.Context, roomID uuid.UUID) ([]domain.Seat, error) {
	results, err := r.q.ListSeatsByRoom(ctx, roomID)
	if err != nil {
		return nil, err
	}
	seats := make([]domain.Seat, len(results))
	for i, row := range results {
		seats[i] = domain.Seat{
			ID:     row.ID,
			RoomID: row.RoomID,
			TeamID: row.TeamID,
			Label:  row.Label,
			Position: domain.Position{
				X: int(row.PosX),
				Y: int(row.PosY),
			},
			Rotation: domain.Rotation(row.Rotation),
		}
	}
	return seats, nil
}

func (r *SeatRepository) Update(ctx context.Context, seat *domain.Seat) error {
	_, err := r.q.UpdateSeat(ctx, sqlc.UpdateSeatParams{
		ID:       seat.ID,
		Label:    seat.Label,
		PosX:     int32(seat.Position.X),
		PosY:     int32(seat.Position.Y),
		Rotation: int32(seat.Rotation),
	})
	return err
}

func (r *SeatRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteSeat(ctx, id)
}
