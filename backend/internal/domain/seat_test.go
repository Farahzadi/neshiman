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
		seat, err := NewSeat(roomID, teamID, "A1", pos, Rotation0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if seat.Label != "A1" {
			t.Errorf("got label %q, want %q", seat.Label, "A1")
		}
		if seat.Position != pos {
			t.Errorf("got position %+v, want %+v", seat.Position, pos)
		}
		if seat.Rotation != Rotation0 {
			t.Errorf("got rotation %d, want %d", seat.Rotation, Rotation0)
		}
		if seat.ID == uuid.Nil {
			t.Error("expected non-nil ID")
		}
	})

	t.Run("valid rotations", func(t *testing.T) {
		rotations := []Rotation{Rotation0, Rotation90, Rotation180, Rotation270}
		for _, rot := range rotations {
			_, err := NewSeat(roomID, teamID, "X", Position{1, 1}, rot)
			if err != nil {
				t.Errorf("rotation %d should be valid, got: %v", rot, err)
			}
		}
	})

	t.Run("invalid rotation", func(t *testing.T) {
		_, err := NewSeat(roomID, teamID, "X", Position{1, 1}, Rotation(45))
		if err != ErrInvalidRotation {
			t.Errorf("got %v, want %v", err, ErrInvalidRotation)
		}
	})
}

func TestSeatMoveTo(t *testing.T) {
	seat, _ := NewSeat(uuid.New(), uuid.New(), "A1", Position{0, 0}, Rotation0)
	seat.MoveTo(Position{X: 10, Y: 20})
	if seat.Position.X != 10 || seat.Position.Y != 20 {
		t.Errorf("got position %+v, want {10 20}", seat.Position)
	}
}

func TestSeatRotate(t *testing.T) {
	seat, _ := NewSeat(uuid.New(), uuid.New(), "A1", Position{0, 0}, Rotation0)
	seat.Rotate(Rotation90)
	if seat.Rotation != Rotation90 {
		t.Errorf("got rotation %d, want %d", seat.Rotation, Rotation90)
	}
}
