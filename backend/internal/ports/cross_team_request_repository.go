package ports

import (
	"context"
	"neshiman/backend/internal/domain"

	"github.com/google/uuid"
)

type CrossTeamRequestRepository interface {
	Create(ctx context.Context, request *domain.CrossTeamRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.CrossTeamRequest, error)
	ListByStatus(ctx context.Context, status domain.RequestStatus) ([]domain.CrossTeamRequest, error)
	ListPendingByTeam(ctx context.Context, teamID uuid.UUID) ([]domain.CrossTeamRequest, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.RequestStatus) error
}
