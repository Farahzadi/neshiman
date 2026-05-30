package ports

import (
	"context"
	"neshiman/backend/internal/domain"

	"github.com/google/uuid"
)

type ReservationWithDetails struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	UserName  string
	SeatID    uuid.UUID
	SeatLabel string
	Date      domain.Date
}

type ReservationRepository interface {
	Create(ctx context.Context, reservation *domain.Reservation) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Reservation, error)
	GetBySeatAndDate(ctx context.Context, seatID uuid.UUID, date domain.Date) (*domain.Reservation, error)
	ListByDate(ctx context.Context, date domain.Date) ([]domain.Reservation, error)
	ListByUserAndDate(ctx context.Context, userID uuid.UUID, date domain.Date) ([]domain.Reservation, error)
	ListByUserAndWeek(ctx context.Context, userID uuid.UUID, start, end domain.Date) ([]domain.Reservation, error)
	CountByUserInWeek(ctx context.Context, userID uuid.UUID, start, end domain.Date) (int, error)
	ListByRoomAndDate(ctx context.Context, roomID uuid.UUID, date domain.Date) ([]domain.Reservation, error)
	ListByRoomAndDateWithDetails(ctx context.Context, roomID uuid.UUID, date domain.Date) ([]ReservationWithDetails, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
