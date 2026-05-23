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

func (r *SeatRepository) BulkSync(ctx context.Context, roomID uuid.UUID, seats []domain.Seat) ([]domain.Seat, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)

	existing, err := q.ListSeatsByRoom(ctx, roomID)
	if err != nil {
		return nil, err
	}

	existingMap := make(map[uuid.UUID]struct{})
	for _, s := range existing {
		existingMap[s.ID] = struct{}{}
	}

	incomingIDs := make(map[uuid.UUID]struct{})
	var result []domain.Seat

	for _, seat := range seats {
		if seat.ID == uuid.Nil {
			created, err := q.CreateSeat(ctx, sqlc.CreateSeatParams{
				RoomID:   roomID,
				TeamID:   seat.TeamID,
				Label:    seat.Label,
				PosX:     int32(seat.Position.X),
				PosY:     int32(seat.Position.Y),
				Rotation: int32(seat.Rotation),
			})
			if err != nil {
				return nil, err
			}
			result = append(result, rowToSeat(created))
		} else {
			incomingIDs[seat.ID] = struct{}{}
			if _, exists := existingMap[seat.ID]; exists {
				updated, err := q.UpdateSeatFull(ctx, sqlc.UpdateSeatFullParams{
					ID:       seat.ID,
					TeamID:   seat.TeamID,
					Label:    seat.Label,
					PosX:     int32(seat.Position.X),
					PosY:     int32(seat.Position.Y),
					Rotation: int32(seat.Rotation),
				})
				if err != nil {
					return nil, err
				}
				result = append(result, rowToSeat(updated))
			}
		}
	}

	for _, s := range existing {
		if _, stillExists := incomingIDs[s.ID]; !stillExists {
			if err := q.DeleteSeat(ctx, s.ID); err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return result, nil
}

func rowToSeat(row sqlc.Seat) domain.Seat {
	return domain.Seat{
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
