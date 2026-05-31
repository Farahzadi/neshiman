package ports

import (
	"context"
	"neshiman/backend/internal/domain"

	"github.com/google/uuid"
)

type SeatRepository interface {
	Create(ctx context.Context, seat *domain.Seat) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Seat, error)
	GetByAssignedUser(ctx context.Context, userID uuid.UUID) (*domain.Seat, error)
	ListByRoom(ctx context.Context, roomID uuid.UUID) ([]domain.Seat, error)
	Update(ctx context.Context, seat *domain.Seat) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByRoom(ctx context.Context, roomID uuid.UUID) error
	BulkSync(ctx context.Context, roomID uuid.UUID, seats []domain.Seat) ([]domain.Seat, error)
	AssignUser(ctx context.Context, seatID, userID uuid.UUID) (*domain.Seat, error)
	UnassignUser(ctx context.Context, seatID uuid.UUID) (*domain.Seat, error)
}
