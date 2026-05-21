package postgres

import (
	"context"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"
	"testing"

	"github.com/google/uuid"
)

func TestTeamRepository(t *testing.T) {
	pool := testPool(t)
	repo := NewTeamRepository(pool)

	ctx := context.Background()

	t.Run("create and get by id", func(t *testing.T) {
		truncate(t, pool)
		team := domain.NewTeam("Engineering")
		if err := repo.Create(ctx, team); err != nil {
			t.Fatalf("create: %v", err)
		}
		if team.ID == uuid.Nil {
			t.Fatal("expected non-nil ID")
		}

		got, err := repo.GetByID(ctx, team.ID)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if got.Name != "Engineering" {
			t.Errorf("got %q, want %q", got.Name, "Engineering")
		}
	})

	t.Run("get not found", func(t *testing.T) {
		truncate(t, pool)
		_, err := repo.GetByID(ctx, uuid.New())
		if err != domain.ErrTeamNotFound {
			t.Errorf("got %v, want %v", err, domain.ErrTeamNotFound)
		}
	})

	t.Run("list", func(t *testing.T) {
		truncate(t, pool)
		repo.Create(ctx, domain.NewTeam("A"))
		repo.Create(ctx, domain.NewTeam("B"))

		teams, err := repo.List(ctx)
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if len(teams) != 2 {
			t.Errorf("got %d, want 2", len(teams))
		}
	})

	t.Run("delete", func(t *testing.T) {
		truncate(t, pool)
		team := domain.NewTeam("ToDelete")
		repo.Create(ctx, team)
		repo.Delete(ctx, team.ID)

		_, err := repo.GetByID(ctx, team.ID)
		if err != domain.ErrTeamNotFound {
			t.Errorf("got %v, want %v", err, domain.ErrTeamNotFound)
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		truncate(t, pool)
		repo.Create(ctx, domain.NewTeam("Unique"))
		err := repo.Create(ctx, domain.NewTeam("Unique"))
		if err == nil {
			t.Error("expected error for duplicate name")
		}
	})
}

func TestTeamRepositoryInterface(t *testing.T) {
	var _ ports.TeamRepository = (*TeamRepository)(nil)
}
