package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewSeat(t *testing.T) {
	roomID := uuid.New()
	teamID := uuid.New()

	t.Run("valid seat", func(t *testing.T) {
		pos := Position{X: 5, Y: 3}
		seat, err := NewSeat(roomID, teamID, "A1", pos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if seat.Label != "A1" {
			t.Errorf("got label %q, want %q", seat.Label, "A1")
		}
		if seat.Position != pos {
			t.Errorf("got position %+v, want %+v", seat.Position, pos)
		}
		if seat.ID == uuid.Nil {
			t.Error("expected non-nil ID")
		}
	})
}

func TestSeatMoveTo(t *testing.T) {
	seat, _ := NewSeat(uuid.New(), uuid.New(), "A1", Position{0, 0})
	seat.MoveTo(Position{X: 10, Y: 20})
	if seat.Position.X != 10 || seat.Position.Y != 20 {
		t.Errorf("got position %+v, want {10 20}", seat.Position)
	}
}
