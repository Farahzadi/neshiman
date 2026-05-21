package application

import (
	"context"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"

	"github.com/google/uuid"
)

type TeamService struct {
	teams ports.TeamRepository
}

func NewTeamService(teams ports.TeamRepository) *TeamService {
	return &TeamService{teams: teams}
}

func (s *TeamService) CreateTeam(ctx context.Context, name string) (*domain.Team, error) {
	team := domain.NewTeam(name)
	if err := s.teams.Create(ctx, team); err != nil {
		return nil, err
	}
	return team, nil
}

func (s *TeamService) GetTeam(ctx context.Context, id uuid.UUID) (*domain.Team, error) {
	return s.teams.GetByID(ctx, id)
}

func (s *TeamService) ListTeams(ctx context.Context) ([]domain.Team, error) {
	return s.teams.List(ctx)
}

func (s *TeamService) DeleteTeam(ctx context.Context, id uuid.UUID) error {
	return s.teams.Delete(ctx, id)
}
