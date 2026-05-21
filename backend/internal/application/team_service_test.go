package application

import (
	"context"
	"errors"
	"neshiman/backend/internal/domain"
	"testing"

	"github.com/google/uuid"
)

func TestTeamService_CreateTeam(t *testing.T) {
	var saved *domain.Team
	svc := NewTeamService(&mockTeamRepo{
		createFn: func(_ context.Context, tm *domain.Team) error {
			saved = tm
			return nil
		},
	})
	team, err := svc.CreateTeam(context.Background(), "Engineering")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if team.Name != "Engineering" {
		t.Errorf("got %q, want %q", team.Name, "Engineering")
	}
	if saved != team {
		t.Error("saved should be same pointer")
	}
}

func TestTeamService_GetTeam(t *testing.T) {
	id := uuid.New()
	svc := NewTeamService(&mockTeamRepo{
		getByIDFn: func(_ context.Context, got uuid.UUID) (*domain.Team, error) {
			if got != id {
				t.Errorf("got %v, want %v", got, id)
			}
			return &domain.Team{ID: id, Name: "Found"}, nil
		},
	})
	team, err := svc.GetTeam(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if team.Name != "Found" {
		t.Errorf("got %q, want %q", team.Name, "Found")
	}
}

func TestTeamService_ListTeams(t *testing.T) {
	svc := NewTeamService(&mockTeamRepo{
		listFn: func(_ context.Context) ([]domain.Team, error) {
			return []domain.Team{{Name: "Eng"}, {Name: "Design"}}, nil
		},
	})
	teams, err := svc.ListTeams(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(teams) != 2 {
		t.Errorf("got %d, want 2", len(teams))
	}
}

func TestTeamService_DeleteTeam(t *testing.T) {
	id := uuid.New()
	var deleted uuid.UUID
	svc := NewTeamService(&mockTeamRepo{
		deleteFn: func(_ context.Context, got uuid.UUID) error {
			deleted = got
			return nil
		},
	})
	if err := svc.DeleteTeam(context.Background(), id); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deleted != id {
		t.Errorf("got %v, want %v", deleted, id)
	}
}

func TestTeamService_RepoError(t *testing.T) {
	svc := NewTeamService(&mockTeamRepo{
		createFn: func(_ context.Context, _ *domain.Team) error {
			return errors.New("db error")
		},
	})
	_, err := svc.CreateTeam(context.Background(), "X")
	if err == nil || err.Error() != "db error" {
		t.Errorf("got %v, want db error", err)
	}
}
