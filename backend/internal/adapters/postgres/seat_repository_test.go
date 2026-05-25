package postgres

import (
	"context"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"
	"testing"

	"github.com/google/uuid"
)

func TestSeatRepository(t *testing.T) {
	pool := testPool(t)
	roomRepo := NewRoomRepository(pool)
	teamRepo := NewTeamRepository(pool)
	seatRepo := NewSeatRepository(pool)

	ctx := context.Background()

	makeRoom := func(t *testing.T) domain.Room {
		t.Helper()
		room, _ := domain.NewRoom("Room-"+uuid.New().String()[:8], 20, 20)
		if err := roomRepo.Create(ctx, room); err != nil {
			t.Fatalf("setup room: %v", err)
		}
		return *room
	}

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
		room := makeRoom(t)
		team := makeTeam(t)

		seat, err := domain.NewSeat(room.ID, team.ID, "A1", domain.Position{X: 5, Y: 3})
		if err != nil {
			t.Fatal(err)
		}
		if err := seatRepo.Create(ctx, seat); err != nil {
			t.Fatalf("create: %v", err)
		}
		if seat.ID == uuid.Nil {
			t.Fatal("expected non-nil ID")
		}

		got, err := seatRepo.GetByID(ctx, seat.ID)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if got.Label != "A1" || got.Position.X != 5 || got.Position.Y != 3 {
			t.Errorf("got %+v, want A1 at (5,3)", got)
		}
	})

	t.Run("get not found", func(t *testing.T) {
		truncate(t, pool)
		_, err := seatRepo.GetByID(ctx, uuid.New())
		if err != domain.ErrSeatNotFound {
			t.Errorf("got %v, want %v", err, domain.ErrSeatNotFound)
		}
	})

	t.Run("list by room", func(t *testing.T) {
		truncate(t, pool)
		room := makeRoom(t)
		team := makeTeam(t)

		seatRepo.Create(ctx, mustNewSeat(t, room.ID, team.ID, "A1", 1, 1))
		seatRepo.Create(ctx, mustNewSeat(t, room.ID, team.ID, "A2", 2, 1))

		seats, err := seatRepo.ListByRoom(ctx, room.ID)
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if len(seats) != 2 {
			t.Errorf("got %d, want 2", len(seats))
		}
	})

	t.Run("update", func(t *testing.T) {
		truncate(t, pool)
		room := makeRoom(t)
		team := makeTeam(t)
		seat := mustNewSeat(t, room.ID, team.ID, "Old", 0, 0)
		seatRepo.Create(ctx, seat)

		seat.Label = "Updated"
		seat.Position = domain.Position{X: 10, Y: 20}
		if err := seatRepo.Update(ctx, seat); err != nil {
			t.Fatalf("update: %v", err)
		}

		got, _ := seatRepo.GetByID(ctx, seat.ID)
		if got.Label != "Updated" || got.Position.X != 10 || got.Position.Y != 20 {
			t.Errorf("got %+v, want Updated at (10,20)", got)
		}
	})

	t.Run("delete", func(t *testing.T) {
		truncate(t, pool)
		room := makeRoom(t)
		team := makeTeam(t)
		seat := mustNewSeat(t, room.ID, team.ID, "Del", 0, 0)
		seatRepo.Create(ctx, seat)
		seatRepo.Delete(ctx, seat.ID)

		_, err := seatRepo.GetByID(ctx, seat.ID)
		if err != domain.ErrSeatNotFound {
			t.Errorf("got %v, want %v", err, domain.ErrSeatNotFound)
		}
	})
}

func mustNewSeat(t *testing.T, roomID, teamID uuid.UUID, label string, x, y int) *domain.Seat {
	t.Helper()
	s, err := domain.NewSeat(roomID, teamID, label, domain.Position{X: x, Y: y})
	if err != nil {
		t.Fatalf("NewSeat: %v", err)
	}
	return s
}

func TestSeatRepository_BulkSync(t *testing.T) {
	pool := testPool(t)
	roomRepo := NewRoomRepository(pool)
	teamRepo := NewTeamRepository(pool)
	seatRepo := NewSeatRepository(pool)

	ctx := context.Background()

	makeRoom := func(t *testing.T) domain.Room {
		t.Helper()
		room, _ := domain.NewRoom("Room-"+uuid.New().String()[:8], 20, 20)
		if err := roomRepo.Create(ctx, room); err != nil {
			t.Fatalf("setup room: %v", err)
		}
		return *room
	}

	makeTeam := func(t *testing.T) domain.Team {
		t.Helper()
		team := domain.NewTeam("Team-" + uuid.New().String()[:8])
		if err := teamRepo.Create(ctx, team); err != nil {
			t.Fatalf("setup team: %v", err)
		}
		return *team
	}

	t.Run("sync seats from empty", func(t *testing.T) {
		truncate(t, pool)
		room := makeRoom(t)
		team := makeTeam(t)

		seats := []domain.Seat{
			{TeamID: team.ID, Label: "S1", Position: domain.Position{X: 0, Y: 0}},
			{TeamID: team.ID, Label: "S2", Position: domain.Position{X: 2, Y: 3}},
		}
		result, err := seatRepo.BulkSync(ctx, room.ID, seats)
		if err != nil {
			t.Fatalf("BulkSync: %v", err)
		}
		if len(result) != 2 {
			t.Fatalf("got %d seats, want 2", len(result))
		}
		for _, s := range result {
			if s.ID == uuid.Nil {
				t.Error("expected non-nil ID on created seat")
			}
		}
	})

	t.Run("sync with updates and deletes", func(t *testing.T) {
		truncate(t, pool)
		room := makeRoom(t)
		team := makeTeam(t)

		existing := mustNewSeat(t, room.ID, team.ID, "Keep", 0, 0)
		toDelete := mustNewSeat(t, room.ID, team.ID, "Del", 5, 5)
		seatRepo.Create(ctx, existing)
		seatRepo.Create(ctx, toDelete)

		updatedTeam := makeTeam(t)
		seats := []domain.Seat{
			{ID: existing.ID, TeamID: updatedTeam.ID, Label: "Changed", Position: domain.Position{X: 10, Y: 10}},
		}
		result, err := seatRepo.BulkSync(ctx, room.ID, seats)
		if err != nil {
			t.Fatalf("BulkSync: %v", err)
		}
		if len(result) != 1 {
			t.Fatalf("got %d seats, want 1", len(result))
		}
		got := result[0]
		if got.ID != existing.ID {
			t.Error("expected same ID for updated seat")
		}
		if got.Label != "Changed" || got.Position.X != 10 || got.Position.Y != 10 {
			t.Errorf("got %+v, want changed seat", got)
		}
		if got.TeamID != updatedTeam.ID {
			t.Errorf("got team %v, want %v", got.TeamID, updatedTeam.ID)
		}

		_, err = seatRepo.GetByID(ctx, toDelete.ID)
		if err != domain.ErrSeatNotFound {
			t.Error("expected deleted seat to be removed")
		}
	})

	t.Run("empty seat list deletes all", func(t *testing.T) {
		truncate(t, pool)
		room := makeRoom(t)
		team := makeTeam(t)

		seatRepo.Create(ctx, mustNewSeat(t, room.ID, team.ID, "A", 0, 0))
		seatRepo.Create(ctx, mustNewSeat(t, room.ID, team.ID, "B", 1, 1))

		result, err := seatRepo.BulkSync(ctx, room.ID, []domain.Seat{})
		if err != nil {
			t.Fatalf("BulkSync: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("got %d seats, want 0", len(result))
		}
	})
}

func TestSeatRepositoryInterface(t *testing.T) {
	var _ ports.SeatRepository = (*SeatRepository)(nil)
}
