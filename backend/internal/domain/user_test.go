package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	tid := uuid.New()

	t.Run("with team", func(t *testing.T) {
		u := NewUser("Alice", "alice@test", &tid, RoleViewer, 2)
		if u.Name != "Alice" || u.Email != "alice@test" {
			t.Errorf("got %s/%s, want Alice/alice@test", u.Name, u.Email)
		}
		if *u.TeamID != tid {
			t.Errorf("got team %v, want %v", u.TeamID, tid)
		}
		if u.Role != RoleViewer {
			t.Errorf("got role %v, want %v", u.Role, RoleViewer)
		}
		if u.WeeklyLimit != 2 {
			t.Errorf("got limit %d, want 2", u.WeeklyLimit)
		}
	})

	t.Run("no team", func(t *testing.T) {
		u := NewUser("Bob", "bob@test", nil, RoleSuperAdmin, 5)
		if u.TeamID != nil {
			t.Error("expected nil team ID")
		}
	})

	t.Run("default weekly limit", func(t *testing.T) {
		u := NewUser("Carol", "carol@test", nil, RoleViewer, 0)
		if u.WeeklyLimit != 2 {
			t.Errorf("got limit %d, want default 2", u.WeeklyLimit)
		}
	})
}

func TestUserIsSuperAdmin(t *testing.T) {
	u := NewUser("Admin", "admin@test", nil, RoleSuperAdmin, 5)
	if !u.IsSuperAdmin() {
		t.Error("expected superadmin")
	}
	if u.IsTeamAdmin() {
		t.Error("superadmin should not be team admin")
	}
}

func TestUserIsTeamAdmin(t *testing.T) {
	u := NewUser("Lead", "lead@test", nil, RoleTeamAdmin, 3)
	if !u.IsTeamAdmin() {
		t.Error("expected team admin")
	}
	if u.IsSuperAdmin() {
		t.Error("team admin should not be superadmin")
	}
}

func TestUserIsAdminOfTeam(t *testing.T) {
	teamA := uuid.New()
	teamB := uuid.New()

	t.Run("superadmin", func(t *testing.T) {
		u := NewUser("Admin", "admin@test", nil, RoleSuperAdmin, 5)
		if !u.IsAdminOfTeam(teamA) {
			t.Error("superadmin should be admin of any team")
		}
	})

	t.Run("team admin owns team", func(t *testing.T) {
		u := NewUser("Lead", "lead@test", &teamA, RoleTeamAdmin, 3)
		if !u.IsAdminOfTeam(teamA) {
			t.Error("team admin should be admin of own team")
		}
	})

	t.Run("team admin other team", func(t *testing.T) {
		u := NewUser("Lead", "lead@test", &teamA, RoleTeamAdmin, 3)
		if u.IsAdminOfTeam(teamB) {
			t.Error("team admin should not be admin of other team")
		}
	})

	t.Run("viewer", func(t *testing.T) {
		u := NewUser("Viewer", "viewer@test", &teamA, RoleViewer, 2)
		if u.IsAdminOfTeam(teamA) {
			t.Error("viewer should not be admin")
		}
	})

	t.Run("team admin nil team", func(t *testing.T) {
		u := NewUser("Lead", "lead@test", nil, RoleTeamAdmin, 3)
		if u.IsAdminOfTeam(teamA) {
			t.Error("admin without team should not be admin of any team")
		}
	})
}
