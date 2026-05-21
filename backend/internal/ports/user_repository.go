package ports

import (
	"context"
	"neshiman/backend/internal/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	ListByTeam(ctx context.Context, teamID uuid.UUID) ([]domain.User, error)
	ListAll(ctx context.Context) ([]domain.User, error)
	UpdateWeeklyLimit(ctx context.Context, userID uuid.UUID, limit domain.WeeklyLimit) error
	Delete(ctx context.Context, id uuid.UUID) error
}
