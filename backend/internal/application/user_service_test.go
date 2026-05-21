package application

import (
	"context"
	"errors"
	"neshiman/backend/internal/domain"
	"testing"

	"github.com/google/uuid"
)

func TestUserService_CreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var saved *domain.User
		svc := NewUserService(&mockUserRepo{
			createFn: func(_ context.Context, u *domain.User) error {
				saved = u
				return nil
			},
		})

		tid := uuid.New()
		user, err := svc.CreateUser(context.Background(), "Alice", "alice@test", &tid, domain.RoleViewer, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.Name != "Alice" || user.Email != "alice@test" {
			t.Error("user fields mismatch")
		}
		if saved != user {
			t.Error("saved user should be the same pointer")
		}
	})

	t.Run("default limit", func(t *testing.T) {
		svc := NewUserService(&mockUserRepo{
			createFn: func(_ context.Context, u *domain.User) error {
				if u.WeeklyLimit != 2 {
					t.Errorf("got limit %d, want default 2", u.WeeklyLimit)
				}
				return nil
			},
		})
		_, err := svc.CreateUser(context.Background(), "Bob", "bob@test", nil, domain.RoleViewer, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		svc := NewUserService(&mockUserRepo{
			createFn: func(_ context.Context, _ *domain.User) error {
				return errors.New("db error")
			},
		})
		_, err := svc.CreateUser(context.Background(), "X", "x@test", nil, domain.RoleViewer, 2)
		if err == nil || err.Error() != "db error" {
			t.Errorf("got %v, want db error", err)
		}
	})
}

func TestUserService_GetUser(t *testing.T) {
	id := uuid.New()
	svc := NewUserService(&mockUserRepo{
		getByIDFn: func(_ context.Context, got uuid.UUID) (*domain.User, error) {
			if got != id {
				t.Errorf("got %v, want %v", got, id)
			}
			return &domain.User{ID: id, Name: "Found"}, nil
		},
	})
	user, err := svc.GetUser(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Name != "Found" {
		t.Errorf("got %q, want %q", user.Name, "Found")
	}
}

func TestUserService_GetUserByEmail(t *testing.T) {
	svc := NewUserService(&mockUserRepo{
		getByEmailFn: func(_ context.Context, email string) (*domain.User, error) {
			if email != "alice@test" {
				t.Errorf("got %q, want alice@test", email)
			}
			return &domain.User{Email: "alice@test"}, nil
		},
	})
	user, err := svc.GetUserByEmail(context.Background(), "alice@test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "alice@test" {
		t.Errorf("got %q, want alice@test", user.Email)
	}
}

func TestUserService_ListUsersByTeam(t *testing.T) {
	teamID := uuid.New()
	svc := NewUserService(&mockUserRepo{
		listByTeamFn: func(_ context.Context, got uuid.UUID) ([]domain.User, error) {
			if got != teamID {
				t.Errorf("got %v, want %v", got, teamID)
			}
			return []domain.User{{Name: "A"}, {Name: "B"}}, nil
		},
	})
	users, err := svc.ListUsersByTeam(context.Background(), teamID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("got %d, want 2", len(users))
	}
}

func TestUserService_UpdateWeeklyLimit(t *testing.T) {
	id := uuid.New()
	svc := NewUserService(&mockUserRepo{
		updateWeeklyLimitFn: func(_ context.Context, got uuid.UUID, limit domain.WeeklyLimit) error {
			if got != id {
				t.Errorf("got %v, want %v", got, id)
			}
			if limit != 5 {
				t.Errorf("got %d, want 5", limit)
			}
			return nil
		},
	})
	if err := svc.UpdateWeeklyLimit(context.Background(), id, 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	id := uuid.New()
	var deleted uuid.UUID
	svc := NewUserService(&mockUserRepo{
		deleteFn: func(_ context.Context, got uuid.UUID) error {
			deleted = got
			return nil
		},
	})
	if err := svc.DeleteUser(context.Background(), id); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deleted != id {
		t.Errorf("got %v, want %v", deleted, id)
	}
}
