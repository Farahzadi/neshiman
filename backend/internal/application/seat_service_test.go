package application

import (
	"context"
	"errors"
	"neshiman/backend/internal/domain"
	"testing"

	"github.com/google/uuid"
)

func TestSeatService_CreateSeat(t *testing.T) {
	roomID := uuid.New()
	teamID := uuid.New()

	t.Run("success", func(t *testing.T) {
		var saved *domain.Seat
		svc := NewSeatService(&mockSeatRepo{
			createFn: func(_ context.Context, s *domain.Seat) error {
				saved = s
				return nil
			},
		})

		seat, err := svc.CreateSeat(context.Background(), roomID, teamID, "A1", 5, 3, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if seat.Label != "A1" || seat.Position.X != 5 || seat.Position.Y != 3 {
			t.Error("seat fields mismatch")
		}
		if saved != seat {
			t.Error("saved seat should be the same pointer")
		}
	})

	t.Run("invalid rotation", func(t *testing.T) {
		svc := NewSeatService(&mockSeatRepo{
			createFn: func(_ context.Context, _ *domain.Seat) error { return nil },
		})
		_, err := svc.CreateSeat(context.Background(), roomID, teamID, "X", 1, 1, 45)
		if !errors.Is(err, domain.ErrInvalidRotation) {
			t.Errorf("got %v, want %v", err, domain.ErrInvalidRotation)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		svc := NewSeatService(&mockSeatRepo{
			createFn: func(_ context.Context, _ *domain.Seat) error {
				return errors.New("db error")
			},
		})
		_, err := svc.CreateSeat(context.Background(), roomID, teamID, "X", 1, 1, 0)
		if err == nil || err.Error() != "db error" {
			t.Errorf("got %v, want db error", err)
		}
	})
}

func TestSeatService_MoveSeat(t *testing.T) {
	id := uuid.New()
	svc := NewSeatService(&mockSeatRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Seat, error) {
			return &domain.Seat{ID: id, Position: domain.Position{X: 0, Y: 0}}, nil
		},
		updateFn: func(_ context.Context, s *domain.Seat) error {
			if s.Position.X != 10 || s.Position.Y != 20 {
				t.Error("position not updated before save")
			}
			return nil
		},
	})

	seat, err := svc.MoveSeat(context.Background(), id, 10, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if seat.Position.X != 10 || seat.Position.Y != 20 {
		t.Errorf("got %+v, want {10 20}", seat.Position)
	}
}

func TestSeatService_RotateSeat(t *testing.T) {
	id := uuid.New()
	svc := NewSeatService(&mockSeatRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Seat, error) {
			return &domain.Seat{ID: id, Rotation: domain.Rotation0}, nil
		},
		updateFn: func(_ context.Context, s *domain.Seat) error {
			if s.Rotation != domain.Rotation180 {
				t.Errorf("got rotation %d, want %d", s.Rotation, domain.Rotation180)
			}
			return nil
		},
	})

	seat, err := svc.RotateSeat(context.Background(), id, 180)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if seat.Rotation != domain.Rotation180 {
		t.Errorf("got %d, want %d", seat.Rotation, domain.Rotation180)
	}
}

func TestSeatService_RotateSeat_Invalid(t *testing.T) {
	id := uuid.New()
	svc := NewSeatService(&mockSeatRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Seat, error) {
			return &domain.Seat{ID: id}, nil
		},
	})
	_, err := svc.RotateSeat(context.Background(), id, 45)
	if !errors.Is(err, domain.ErrInvalidRotation) {
		t.Errorf("got %v, want %v", err, domain.ErrInvalidRotation)
	}
}

func TestSeatService_GetSeat(t *testing.T) {
	id := uuid.New()
	svc := NewSeatService(&mockSeatRepo{
		getByIDFn: func(_ context.Context, got uuid.UUID) (*domain.Seat, error) {
			if got != id {
				t.Errorf("got id %v, want %v", got, id)
			}
			return &domain.Seat{ID: id, Label: "Found"}, nil
		},
	})

	seat, err := svc.GetSeat(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if seat.Label != "Found" {
		t.Errorf("got %q, want %q", seat.Label, "Found")
	}
}

func TestSeatService_ListSeatsByRoom(t *testing.T) {
	roomID := uuid.New()
	svc := NewSeatService(&mockSeatRepo{
		listByRoomFn: func(_ context.Context, got uuid.UUID) ([]domain.Seat, error) {
			if got != roomID {
				t.Errorf("got %v, want %v", got, roomID)
			}
			return []domain.Seat{{Label: "A1"}, {Label: "A2"}}, nil
		},
	})

	seats, err := svc.ListSeatsByRoom(context.Background(), roomID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(seats) != 2 {
		t.Errorf("got %d seats, want 2", len(seats))
	}
}

func TestSeatService_DeleteSeat(t *testing.T) {
	id := uuid.New()
	var deleted uuid.UUID
	svc := NewSeatService(&mockSeatRepo{
		deleteFn: func(_ context.Context, got uuid.UUID) error {
			deleted = got
			return nil
		},
	})

	if err := svc.DeleteSeat(context.Background(), id); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deleted != id {
		t.Errorf("got %v, want %v", deleted, id)
	}
}
