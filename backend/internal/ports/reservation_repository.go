package ports

import (
	"context"
	"neshiman/backend/internal/domain"

	"github.com/google/uuid"
)

type ReservationRepository interface {
	Create(ctx context.Context, reservation *domain.Reservation) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Reservation, error)
	GetBySeatAndDate(ctx context.Context, seatID uuid.UUID, date domain.Date) (*domain.Reservation, error)
	ListByUserAndWeek(ctx context.Context, userID uuid.UUID, start, end domain.Date) ([]domain.Reservation, error)
	CountByUserInWeek(ctx context.Context, userID uuid.UUID, start, end domain.Date) (int, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
