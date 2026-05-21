package ports

import (
	"context"
	"neshiman/backend/internal/domain"

	"github.com/google/uuid"
)

type RoomRepository interface {
	Create(ctx context.Context, room *domain.Room) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Room, error)
	List(ctx context.Context) ([]domain.Room, error)
	Update(ctx context.Context, room *domain.Room) error
	Delete(ctx context.Context, id uuid.UUID) error
}
