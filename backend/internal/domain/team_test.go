package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewTeam(t *testing.T) {
	team := NewTeam("Engineering")
	if team.Name != "Engineering" {
		t.Errorf("got name %q, want %q", team.Name, "Engineering")
	}
	if team.ID == uuid.Nil {
		t.Error("expected non-nil ID")
	}
}
