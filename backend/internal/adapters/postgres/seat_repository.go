package postgres

import (
	"context"
	"neshiman/backend/db/sqlc"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
	var assignedUserID pgtype.UUID
	if seat.AssignedUserID != nil {
		assignedUserID = pgtype.UUID{Bytes: *seat.AssignedUserID, Valid: true}
	}
	result, err := r.q.CreateSeat(ctx, sqlc.CreateSeatParams{
		RoomID:         seat.RoomID,
		TeamID:         seat.TeamID,
		Label:          seat.Label,
		PosX:           int32(seat.Position.X),
		PosY:           int32(seat.Position.Y),
		AssignedUserID: assignedUserID,
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
	return rowToSeat(result), nil
}

func (r *SeatRepository) GetByAssignedUser(ctx context.Context, userID uuid.UUID) (*domain.Seat, error) {
	result, err := r.q.GetSeatByAssignedUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrSeatNotFound
		}
		return nil, err
	}
	return rowToSeat(result), nil
}

func (r *SeatRepository) ListByRoom(ctx context.Context, roomID uuid.UUID) ([]domain.Seat, error) {
	results, err := r.q.ListSeatsByRoom(ctx, roomID)
	if err != nil {
		return nil, err
	}
	seats := make([]domain.Seat, len(results))
	for i, row := range results {
		seats[i] = *rowToSeat(row)
	}
	return seats, nil
}

func (r *SeatRepository) Update(ctx context.Context, seat *domain.Seat) error {
	_, err := r.q.UpdateSeat(ctx, sqlc.UpdateSeatParams{
		ID:    seat.ID,
		Label: seat.Label,
		PosX:  int32(seat.Position.X),
		PosY:  int32(seat.Position.Y),
	})
	return err
}

func (r *SeatRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.SoftDeleteSeat(ctx, id)
}

func (r *SeatRepository) DeleteByRoom(ctx context.Context, roomID uuid.UUID) error {
	return r.q.SoftDeleteSeatsByRoom(ctx, roomID)
}

func (r *SeatRepository) AssignUser(ctx context.Context, seatID, userID uuid.UUID) (*domain.Seat, error) {
	result, err := r.q.AssignSeat(ctx, sqlc.AssignSeatParams{
		ID:             seatID,
		AssignedUserID: pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrSeatNotFound
		}
		return nil, err
	}
	return rowToSeat(result), nil
}

func (r *SeatRepository) UnassignUser(ctx context.Context, seatID uuid.UUID) (*domain.Seat, error) {
	result, err := r.q.UnassignSeat(ctx, seatID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrSeatNotFound
		}
		return nil, err
	}
	return rowToSeat(result), nil
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
				RoomID:         roomID,
				TeamID:         seat.TeamID,
				Label:          seat.Label,
				PosX:           int32(seat.Position.X),
				PosY:           int32(seat.Position.Y),
				AssignedUserID: pgtype.UUID{},
			})
			if err != nil {
				return nil, err
			}
			result = append(result, *rowToSeat(created))
		} else {
			incomingIDs[seat.ID] = struct{}{}
			if _, exists := existingMap[seat.ID]; exists {
				updated, err := q.UpdateSeatFull(ctx, sqlc.UpdateSeatFullParams{
					ID:     seat.ID,
					TeamID: seat.TeamID,
					Label:  seat.Label,
					PosX:   int32(seat.Position.X),
					PosY:   int32(seat.Position.Y),
				})
				if err != nil {
					return nil, err
				}
				result = append(result, *rowToSeat(updated))
			}
		}
	}

	for _, s := range existing {
		if _, stillExists := incomingIDs[s.ID]; !stillExists {
			if err := q.SoftDeleteSeat(ctx, s.ID); err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return result, nil
}

func rowToSeat(row sqlc.Seat) *domain.Seat {
	var assignedUserID *uuid.UUID
	if row.AssignedUserID.Valid {
		id := uuid.UUID(row.AssignedUserID.Bytes)
		assignedUserID = &id
	}
	return &domain.Seat{
		ID:     row.ID,
		RoomID: row.RoomID,
		TeamID: row.TeamID,
		Label:  row.Label,
		Position: domain.Position{
			X: int(row.PosX),
			Y: int(row.PosY),
		},
		AssignedUserID: assignedUserID,
	}
}
