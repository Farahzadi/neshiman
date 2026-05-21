package postgres

import (
	"context"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"
	"testing"

	"github.com/google/uuid"
)

func TestUserRepository(t *testing.T) {
	pool := testPool(t)

	// Need a team for FK references
	teamRepo := NewTeamRepository(pool)
	userRepo := NewUserRepository(pool)

	ctx := context.Background()

	makeTeam := func(t *testing.T) domain.Team {
		t.Helper()
		team := domain.NewTeam("Team-" + uuid.New().String()[:8])
		if err := teamRepo.Create(ctx, team); err != nil {
			t.Fatalf("setup team: %v", err)
		}
		return *team
	}

	t.Run("create and get by id", func(t *testing.T) {
		truncate(t, pool)
		team := makeTeam(t)
		user := domain.NewUser("Alice", "alice@test", &team.ID, domain.RoleViewer, 2)

		if err := userRepo.Create(ctx, user); err != nil {
			t.Fatalf("create: %v", err)
		}
		if user.ID == uuid.Nil {
			t.Fatal("expected non-nil ID")
		}

		got, err := userRepo.GetByID(ctx, user.ID)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if got.Name != "Alice" || got.Email != "alice@test" {
			t.Errorf("got %s/%s, want Alice/alice@test", got.Name, got.Email)
		}
		if *got.TeamID != team.ID {
			t.Errorf("got team %v, want %v", got.TeamID, team.ID)
		}
		if got.Role != domain.RoleViewer {
			t.Errorf("got role %v, want viewer", got.Role)
		}
	})

	t.Run("create without team", func(t *testing.T) {
		truncate(t, pool)
		user := domain.NewUser("Admin", "admin@test", nil, domain.RoleSuperAdmin, 5)
		if err := userRepo.Create(ctx, user); err != nil {
			t.Fatalf("create: %v", err)
		}
		got, _ := userRepo.GetByID(ctx, user.ID)
		if got.TeamID != nil {
			t.Error("expected nil team ID")
		}
	})

	t.Run("get by email", func(t *testing.T) {
		truncate(t, pool)
		user := domain.NewUser("Bob", "bob@test", nil, domain.RoleViewer, 2)
		userRepo.Create(ctx, user)

		got, err := userRepo.GetByEmail(ctx, "bob@test")
		if err != nil {
			t.Fatalf("get by email: %v", err)
		}
		if got.Email != "bob@test" {
			t.Errorf("got %q, want bob@test", got.Email)
		}
	})

	t.Run("get by email not found", func(t *testing.T) {
		truncate(t, pool)
		_, err := userRepo.GetByEmail(ctx, "nonexistent@test")
		if err != domain.ErrUserNotFound {
			t.Errorf("got %v, want %v", err, domain.ErrUserNotFound)
		}
	})

	t.Run("list by team", func(t *testing.T) {
		truncate(t, pool)
		teamA := makeTeam(t)
		teamB := makeTeam(t)

		userRepo.Create(ctx, domain.NewUser("A1", "a1@test", &teamA.ID, domain.RoleViewer, 2))
		userRepo.Create(ctx, domain.NewUser("A2", "a2@test", &teamA.ID, domain.RoleViewer, 2))
		userRepo.Create(ctx, domain.NewUser("B1", "b1@test", &teamB.ID, domain.RoleViewer, 2))

		users, err := userRepo.ListByTeam(ctx, teamA.ID)
		if err != nil {
			t.Fatalf("list by team: %v", err)
		}
		if len(users) != 2 {
			t.Errorf("got %d, want 2", len(users))
		}
	})

	t.Run("update weekly limit", func(t *testing.T) {
		truncate(t, pool)
		user := domain.NewUser("Carol", "carol@test", nil, domain.RoleViewer, 2)
		userRepo.Create(ctx, user)

		if err := userRepo.UpdateWeeklyLimit(ctx, user.ID, 5); err != nil {
			t.Fatalf("update limit: %v", err)
		}

		got, _ := userRepo.GetByID(ctx, user.ID)
		if got.WeeklyLimit != 5 {
			t.Errorf("got limit %d, want 5", got.WeeklyLimit)
		}
	})

	t.Run("delete", func(t *testing.T) {
		truncate(t, pool)
		user := domain.NewUser("Del", "del@test", nil, domain.RoleViewer, 2)
		userRepo.Create(ctx, user)
		userRepo.Delete(ctx, user.ID)

		_, err := userRepo.GetByID(ctx, user.ID)
		if err != domain.ErrUserNotFound {
			t.Errorf("got %v, want %v", err, domain.ErrUserNotFound)
		}
	})
}

func TestUserRepositoryInterface(t *testing.T) {
	var _ ports.UserRepository = (*UserRepository)(nil)
}
