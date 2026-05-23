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

type TeamRepository struct {
	q    *sqlc.Queries
	pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) ports.TeamRepository {
	return &TeamRepository{
		q:    sqlc.New(pool),
		pool: pool,
	}
}

func (r *TeamRepository) Create(ctx context.Context, team *domain.Team) error {
	result, err := r.q.CreateTeam(ctx, team.Name)
	if err != nil {
		return err
	}
	team.ID = result.ID
	return nil
}

func (r *TeamRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Team, error) {
	result, err := r.q.GetTeamByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrTeamNotFound
		}
		return nil, err
	}
	return &domain.Team{
		ID:   result.ID,
		Name: result.Name,
	}, nil
}

func (r *TeamRepository) List(ctx context.Context) ([]domain.Team, error) {
	results, err := r.q.ListTeams(ctx)
	if err != nil {
		return nil, err
	}
	teams := make([]domain.Team, len(results))
	for i, row := range results {
		teams[i] = domain.Team{
			ID:   row.ID,
			Name: row.Name,
		}
	}
	return teams, nil
}

func (r *TeamRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.SoftDeleteTeam(ctx, id)
}
