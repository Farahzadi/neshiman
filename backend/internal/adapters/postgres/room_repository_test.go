package postgres

import (
	"context"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"
	"testing"

	"github.com/google/uuid"
)

func TestRoomRepository(t *testing.T) {
	pool := testPool(t)
	repo := NewRoomRepository(pool)

	ctx := context.Background()

	t.Run("create and get by id", func(t *testing.T) {
		truncate(t, pool)
		room, err := domain.NewRoom("Test Room", 10, 8)
		if err != nil {
			t.Fatal(err)
		}
		if err := repo.Create(ctx, room); err != nil {
			t.Fatalf("create: %v", err)
		}
		if room.ID == uuid.Nil {
			t.Fatal("expected non-nil ID after create")
		}

		got, err := repo.GetByID(ctx, room.ID)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if got.Name != "Test Room" || got.GridWidth != 10 || got.GridHeight != 8 {
			t.Errorf("got %+v, want name=Test Room 10x8", got)
		}
	})

	t.Run("get not found", func(t *testing.T) {
		truncate(t, pool)
		_, err := repo.GetByID(ctx, uuid.New())
		if err != domain.ErrRoomNotFound {
			t.Errorf("got %v, want %v", err, domain.ErrRoomNotFound)
		}
	})

	t.Run("list", func(t *testing.T) {
		truncate(t, pool)
		room1, _ := domain.NewRoom("A", 5, 5)
		room2, _ := domain.NewRoom("B", 3, 3)
		repo.Create(ctx, room1)
		repo.Create(ctx, room2)

		rooms, err := repo.List(ctx)
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if len(rooms) != 2 {
			t.Errorf("got %d rooms, want 2", len(rooms))
		}
	})

	t.Run("update", func(t *testing.T) {
		truncate(t, pool)
		room, _ := domain.NewRoom("Old", 5, 5)
		repo.Create(ctx, room)

		room.Name = "Updated"
		room.GridWidth = 20
		room.GridHeight = 15
		if err := repo.Update(ctx, room); err != nil {
			t.Fatalf("update: %v", err)
		}

		got, _ := repo.GetByID(ctx, room.ID)
		if got.Name != "Updated" || got.GridWidth != 20 || got.GridHeight != 15 {
			t.Errorf("got %+v, want Updated 20x15", got)
		}
	})

	t.Run("delete", func(t *testing.T) {
		truncate(t, pool)
		room, _ := domain.NewRoom("ToDelete", 5, 5)
		repo.Create(ctx, room)

		if err := repo.Delete(ctx, room.ID); err != nil {
			t.Fatalf("delete: %v", err)
		}

		_, err := repo.GetByID(ctx, room.ID)
		if err != domain.ErrRoomNotFound {
			t.Errorf("got %v, want %v", err, domain.ErrRoomNotFound)
		}
	})
}

func TestRoomRepositoryInterface(t *testing.T) {
	var _ ports.RoomRepository = (*RoomRepository)(nil)
}
