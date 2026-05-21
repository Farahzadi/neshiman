package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewRoom(t *testing.T) {
	t.Run("valid room", func(t *testing.T) {
		room, err := NewRoom("Test Room", 10, 8)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if room.Name != "Test Room" {
			t.Errorf("got name %q, want %q", room.Name, "Test Room")
		}
		if room.GridWidth != 10 || room.GridHeight != 8 {
			t.Errorf("got grid %dx%d, want 10x8", room.GridWidth, room.GridHeight)
		}
		if room.ID == uuid.Nil {
			t.Error("expected non-nil ID")
		}
	})

	t.Run("zero grid", func(t *testing.T) {
		_, err := NewRoom("Bad", 0, 10)
		if err != ErrInvalidGridSize {
			t.Errorf("got %v, want %v", err, ErrInvalidGridSize)
		}
	})

	t.Run("negative grid", func(t *testing.T) {
		_, err := NewRoom("Bad", 10, -1)
		if err != ErrInvalidGridSize {
			t.Errorf("got %v, want %v", err, ErrInvalidGridSize)
		}
	})

	t.Run("both zero", func(t *testing.T) {
		_, err := NewRoom("Bad", 0, 0)
		if err != ErrInvalidGridSize {
			t.Errorf("got %v, want %v", err, ErrInvalidGridSize)
		}
	})
}
