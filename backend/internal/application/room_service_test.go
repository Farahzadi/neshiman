package application

import (
	"context"
	"errors"
	"neshiman/backend/internal/domain"
	"testing"

	"github.com/google/uuid"
)

func defaultDeleteByRoomFn(_ context.Context, _ uuid.UUID) error { return nil }

func TestRoomService_CreateRoom(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var savedRoom *domain.Room
		svc := NewRoomService(&mockRoomRepo{
			createFn: func(_ context.Context, r *domain.Room) error {
				savedRoom = r
				return nil
			},
		}, &mockSeatRepo{deleteByRoomFn: defaultDeleteByRoomFn})

		room, err := svc.CreateRoom(context.Background(), "Test Room", 10, 8)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if room.Name != "Test Room" || room.GridWidth != 10 || room.GridHeight != 8 {
			t.Error("room fields mismatch")
		}
		if savedRoom != room {
			t.Error("saved room should be the same pointer")
		}
	})

	t.Run("invalid grid", func(t *testing.T) {
		svc := NewRoomService(&mockRoomRepo{createFn: func(_ context.Context, _ *domain.Room) error { return nil }}, &mockSeatRepo{deleteByRoomFn: defaultDeleteByRoomFn})
		_, err := svc.CreateRoom(context.Background(), "Bad", 0, 10)
		if !errors.Is(err, domain.ErrInvalidGridSize) {
			t.Errorf("got %v, want %v", err, domain.ErrInvalidGridSize)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		svc := NewRoomService(&mockRoomRepo{
			createFn: func(_ context.Context, _ *domain.Room) error {
				return errors.New("db error")
			},
		}, &mockSeatRepo{deleteByRoomFn: defaultDeleteByRoomFn})
		_, err := svc.CreateRoom(context.Background(), "Test", 10, 8)
		if err == nil || err.Error() != "db error" {
			t.Errorf("got %v, want db error", err)
		}
	})
}

func TestRoomService_GetRoom(t *testing.T) {
	id := uuid.New()
	svc := NewRoomService(&mockRoomRepo{
		getByIDFn: func(_ context.Context, got uuid.UUID) (*domain.Room, error) {
			if got != id {
				t.Errorf("got id %v, want %v", got, id)
			}
			return &domain.Room{ID: id, Name: "Found"}, nil
		},
	}, &mockSeatRepo{deleteByRoomFn: defaultDeleteByRoomFn})

	room, err := svc.GetRoom(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if room.Name != "Found" {
		t.Errorf("got %q, want %q", room.Name, "Found")
	}
}

func TestRoomService_ListRooms(t *testing.T) {
	svc := NewRoomService(&mockRoomRepo{
		listFn: func(_ context.Context) ([]domain.Room, error) {
			return []domain.Room{{Name: "A"}, {Name: "B"}}, nil
		},
	}, &mockSeatRepo{deleteByRoomFn: defaultDeleteByRoomFn})

	rooms, err := svc.ListRooms(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rooms) != 2 {
		t.Errorf("got %d rooms, want 2", len(rooms))
	}
}

func TestRoomService_UpdateRoom(t *testing.T) {
	id := uuid.New()
	svc := NewRoomService(&mockRoomRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Room, error) {
			return &domain.Room{ID: id, Name: "Old", GridWidth: 5, GridHeight: 5}, nil
		},
		updateFn: func(_ context.Context, r *domain.Room) error {
			if r.Name != "New" || r.GridWidth != 10 || r.GridHeight != 8 {
				t.Error("updated room fields mismatch")
			}
			return nil
		},
	}, &mockSeatRepo{deleteByRoomFn: defaultDeleteByRoomFn})

	room, err := svc.UpdateRoom(context.Background(), id, "New", 10, 8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if room.Name != "New" {
		t.Errorf("got %q, want %q", room.Name, "New")
	}
}

func TestRoomService_DeleteRoom(t *testing.T) {
	id := uuid.New()
	var deleted uuid.UUID
	var deletedByRoom uuid.UUID

	svc := NewRoomService(&mockRoomRepo{
		deleteFn: func(_ context.Context, got uuid.UUID) error {
			deleted = got
			return nil
		},
	}, &mockSeatRepo{
		listByRoomFn: func(_ context.Context, got uuid.UUID) ([]domain.Seat, error) {
			return []domain.Seat{
				{ID: uuid.New(), Label: "S1"},
				{ID: uuid.New(), Label: "S2"},
			}, nil
		},
		deleteByRoomFn: func(_ context.Context, got uuid.UUID) error {
			deletedByRoom = got
			return nil
		},
	})

	n, err := svc.DeleteRoom(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("got %d seats deleted, want 2", n)
	}
	if deletedByRoom != id {
		t.Errorf("deleteByRoom id %v, want %v", deletedByRoom, id)
	}
	if deleted != id {
		t.Errorf("delete id %v, want %v", deleted, id)
	}
}
