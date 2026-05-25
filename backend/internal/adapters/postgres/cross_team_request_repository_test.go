package postgres

import (
	"context"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"
	"testing"

	"github.com/google/uuid"
)

func TestCrossTeamRequestRepository(t *testing.T) {
	pool := testPool(t)

	roomRepo := NewRoomRepository(pool)
	teamRepo := NewTeamRepository(pool)
	seatRepo := NewSeatRepository(pool)
	userRepo := NewUserRepository(pool)
	ctrRepo := NewCrossTeamRequestRepository(pool)

	ctx := context.Background()

	setup := func(t *testing.T, suffix string) (domain.Seat, domain.User) {
		t.Helper()
		room, _ := domain.NewRoom("CTRTestRoom-"+suffix, 10, 10)
		if err := roomRepo.Create(ctx, room); err != nil {
			t.Fatalf("setup room: %v", err)
		}

		team := domain.NewTeam("CTRTestTeam-" + suffix)
		if err := teamRepo.Create(ctx, team); err != nil {
			t.Fatalf("setup team: %v", err)
		}

		seat, _ := domain.NewSeat(room.ID, team.ID, "CTR1-"+suffix, domain.Position{X: 1, Y: 1})
		if err := seatRepo.Create(ctx, seat); err != nil {
			t.Fatalf("setup seat: %v", err)
		}

		user := domain.NewUser("CTRUser-"+suffix, "ctr-"+suffix+"@test", &team.ID, domain.RoleViewer, 2)
		if err := userRepo.Create(ctx, user); err != nil {
			t.Fatalf("setup user: %v", err)
		}

		return *seat, *user
	}

	t.Run("create and get by id", func(t *testing.T) {
		truncate(t, pool)
		seat, user := setup(t, "c1")
		date := domain.Date{Year: 2026, Month: 6, Day: 15}

		req := domain.NewCrossTeamRequest(user.ID, seat.ID, date)
		if err := ctrRepo.Create(ctx, req); err != nil {
			t.Fatalf("create: %v", err)
		}
		if req.ID == uuid.Nil {
			t.Fatal("expected non-nil ID")
		}
		if req.Status != domain.RequestPending {
			t.Errorf("got %v, want pending", req.Status)
		}

		got, err := ctrRepo.GetByID(ctx, req.ID)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if got.RequestingUserID != user.ID || got.TargetSeatID != seat.ID {
			t.Errorf("got requesting_user=%v target_seat=%v, want %v/%v",
				got.RequestingUserID, got.TargetSeatID, user.ID, seat.ID)
		}
	})

	t.Run("get not found", func(t *testing.T) {
		truncate(t, pool)
		_, err := ctrRepo.GetByID(ctx, uuid.New())
		if err != domain.ErrCrossTeamRequestNotFound {
			t.Errorf("got %v, want %v", err, domain.ErrCrossTeamRequestNotFound)
		}
	})

	t.Run("list by status", func(t *testing.T) {
		truncate(t, pool)
		seat, user := setup(t, "lbs1")
		date := domain.Date{Year: 2026, Month: 6, Day: 16}

		req1 := domain.NewCrossTeamRequest(user.ID, seat.ID, date)
		ctrRepo.Create(ctx, req1)


		// Create another request with different date
		seat2, user2 := setup(t, "lbs2")
		req2 := domain.NewCrossTeamRequest(user2.ID, seat2.ID, domain.Date{Year: 2026, Month: 6, Day: 17})
		ctrRepo.Create(ctx, req2)
		ctrRepo.UpdateStatus(ctx, req2.ID, domain.RequestApproved)

		pending, err := ctrRepo.ListByStatus(ctx, domain.RequestPending)
		if err != nil {
			t.Fatalf("list pending: %v", err)
		}
		if len(pending) != 1 {
			t.Errorf("got %d pending, want 1", len(pending))
		}

		approved, err := ctrRepo.ListByStatus(ctx, domain.RequestApproved)
		if err != nil {
			t.Fatalf("list approved: %v", err)
		}
		if len(approved) != 1 {
			t.Errorf("got %d approved, want 1", len(approved))
		}
	})

	t.Run("list pending by team", func(t *testing.T) {
		truncate(t, pool)
		seat, user := setup(t, "lpbt1")
		date := domain.Date{Year: 2026, Month: 6, Day: 18}

		req := domain.NewCrossTeamRequest(user.ID, seat.ID, date)
		ctrRepo.Create(ctx, req)

		// Get the team ID from the seat
		gotSeat, _ := seatRepo.GetByID(ctx, seat.ID)

		requests, err := ctrRepo.ListPendingByTeam(ctx, gotSeat.TeamID)
		if err != nil {
			t.Fatalf("list pending by team: %v", err)
		}
		if len(requests) != 1 {
			t.Errorf("got %d, want 1", len(requests))
		}
		if requests[0].ID != req.ID {
			t.Errorf("got %v, want %v", requests[0].ID, req.ID)
		}
	})

	t.Run("update status", func(t *testing.T) {
		truncate(t, pool)
		seat, user := setup(t, "us1")
		req := domain.NewCrossTeamRequest(user.ID, seat.ID, domain.Date{Year: 2026, Month: 6, Day: 19})
		ctrRepo.Create(ctx, req)

		if err := ctrRepo.UpdateStatus(ctx, req.ID, domain.RequestApproved); err != nil {
			t.Fatalf("update status: %v", err)
		}

		got, _ := ctrRepo.GetByID(ctx, req.ID)
		if got.Status != domain.RequestApproved {
			t.Errorf("got %v, want %v", got.Status, domain.RequestApproved)
		}
	})
}

func TestCrossTeamRequestRepositoryInterface(t *testing.T) {
	var _ ports.CrossTeamRequestRepository = (*CrossTeamRequestRepository)(nil)
}
