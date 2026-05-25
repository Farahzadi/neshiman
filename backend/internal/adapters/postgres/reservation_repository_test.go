package postgres

import (
	"context"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestReservationRepository(t *testing.T) {
	pool := testPool(t)
	roomRepo := NewRoomRepository(pool)
	teamRepo := NewTeamRepository(pool)
	seatRepo := NewSeatRepository(pool)
	userRepo := NewUserRepository(pool)
	reservationRepo := NewReservationRepository(pool)

	ctx := context.Background()

	setup := func(t *testing.T, suffix string) (domain.Room, domain.Team, domain.Seat, domain.User) {
		t.Helper()
		room, _ := domain.NewRoom("ResTestRoom-"+suffix, 10, 10)
		if err := roomRepo.Create(ctx, room); err != nil {
			t.Fatalf("setup room: %v", err)
		}

		team := domain.NewTeam("ResTestTeam-" + suffix)
		if err := teamRepo.Create(ctx, team); err != nil {
			t.Fatalf("setup team: %v", err)
		}

		seat, _ := domain.NewSeat(room.ID, team.ID, "R1-"+suffix, domain.Position{X: 1, Y: 1})
		if err := seatRepo.Create(ctx, seat); err != nil {
			t.Fatalf("setup seat: %v", err)
		}

		user := domain.NewUser("ResUser-"+suffix, "res-"+suffix+"@test", &team.ID, domain.RoleViewer, 5)
		if err := userRepo.Create(ctx, user); err != nil {
			t.Fatalf("setup user: %v", err)
		}

		return *room, *team, *seat, *user
	}

	makeDate := func(year, month, day int) domain.Date {
		return domain.Date{Year: year, Month: month, Day: day}
	}

	t.Run("create and get by id", func(t *testing.T) {
		truncate(t, pool)
		_, _, seat, user := setup(t, "c1")
		date := makeDate(2026, 5, 18)

		reservation := domain.NewReservation(user.ID, seat.ID, date)
		if err := reservationRepo.Create(ctx, reservation); err != nil {
			t.Fatalf("create: %v", err)
		}
		if reservation.ID == uuid.Nil {
			t.Fatal("expected non-nil ID")
		}

		got, err := reservationRepo.GetByID(ctx, reservation.ID)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if got.UserID != user.ID || got.SeatID != seat.ID || got.Date != date {
			t.Errorf("got %+v, want user=%v seat=%v date=%+v", got, user.ID, seat.ID, date)
		}
	})

	t.Run("get not found", func(t *testing.T) {
		truncate(t, pool)
		_, err := reservationRepo.GetByID(ctx, uuid.New())
		if err != domain.ErrReservationNotFound {
			t.Errorf("got %v, want %v", err, domain.ErrReservationNotFound)
		}
	})

	t.Run("get by seat and date", func(t *testing.T) {
		truncate(t, pool)
		_, _, seat, user := setup(t, "gbd1")
		date := makeDate(2026, 5, 15)

		reservation := domain.NewReservation(user.ID, seat.ID, date)
		reservationRepo.Create(ctx, reservation)

		got, err := reservationRepo.GetBySeatAndDate(ctx, seat.ID, date)
		if err != nil {
			t.Fatalf("get by seat and date: %v", err)
		}
		if got == nil {
			t.Fatal("expected non-nil reservation")
		}
		if got.ID != reservation.ID {
			t.Errorf("got %v, want %v", got.ID, reservation.ID)
		}
	})

	t.Run("get by seat and date not found", func(t *testing.T) {
		truncate(t, pool)
		_, _, seat, user := setup(t, "gbdnf1")
		reservation := domain.NewReservation(user.ID, seat.ID, makeDate(2026, 5, 15))
		reservationRepo.Create(ctx, reservation)

		got, err := reservationRepo.GetBySeatAndDate(ctx, seat.ID, makeDate(2026, 5, 16))
		if err != nil {
			t.Fatalf("get by seat and date: %v", err)
		}
		if got != nil {
			t.Error("expected nil for different date")
		}
	})

	t.Run("unique constraint seat+date", func(t *testing.T) {
		truncate(t, pool)
		_, _, seat, user := setup(t, "uc1")
		date := makeDate(2026, 5, 15)
		reservationRepo.Create(ctx, domain.NewReservation(user.ID, seat.ID, date))

		err := reservationRepo.Create(ctx, domain.NewReservation(user.ID, seat.ID, date))
		if err == nil {
			t.Error("expected error for duplicate seat+date")
		}
	})

	t.Run("list by user and week", func(t *testing.T) {
		truncate(t, pool)
		_, _, seatA, user := setup(t, "luw1")

		// Create another seat in same room
		room2, team2, _, _ := setup(t, "luw2")
		seatB, _ := domain.NewSeat(room2.ID, team2.ID, "R2", domain.Position{X: 2, Y: 1})
		seatRepo.Create(ctx, seatB)

		// Mon May 18 2026
		monday := makeDate(2026, 5, 18)
		tuesday := makeDate(2026, 5, 19)
		nextMonday := makeDate(2026, 5, 25)

		reservationRepo.Create(ctx, domain.NewReservation(user.ID, seatA.ID, monday))
		reservationRepo.Create(ctx, domain.NewReservation(user.ID, seatB.ID, tuesday))
		reservationRepo.Create(ctx, domain.NewReservation(user.ID, seatA.ID, nextMonday))

		// List reservations for the week of May 18 (Mon-Sun)
		results, err := reservationRepo.ListByUserAndWeek(ctx, user.ID, monday, makeDate(2026, 5, 24))
		if err != nil {
			t.Fatalf("list by user and week: %v", err)
		}
		if len(results) != 2 {
			t.Errorf("got %d reservations, want 2", len(results))
		}
	})

	t.Run("count by user in week", func(t *testing.T) {
		truncate(t, pool)
		_, _, seatA, user := setup(t, "cwu1")
		room3, team3, _, _ := setup(t, "cwu2")
		seatB, _ := domain.NewSeat(room3.ID, team3.ID, "R2", domain.Position{X: 2, Y: 1})
		seatRepo.Create(ctx, seatB)

		monday := makeDate(2026, 5, 18)
		tuesday := makeDate(2026, 5, 19)

		reservationRepo.Create(ctx, domain.NewReservation(user.ID, seatA.ID, monday))
		reservationRepo.Create(ctx, domain.NewReservation(user.ID, seatB.ID, tuesday))

		count, err := reservationRepo.CountByUserInWeek(ctx, user.ID, monday, makeDate(2026, 5, 24))
		if err != nil {
			t.Fatalf("count: %v", err)
		}
		if count != 2 {
			t.Errorf("got %d, want 2", count)
		}
	})

	t.Run("delete", func(t *testing.T) {
		truncate(t, pool)
		_, _, seat, user := setup(t, "del1")
		date := makeDate(2026, 5, 18)
		reservation := domain.NewReservation(user.ID, seat.ID, date)
		reservationRepo.Create(ctx, reservation)

		if err := reservationRepo.Delete(ctx, reservation.ID); err != nil {
			t.Fatalf("delete: %v", err)
		}

		_, err := reservationRepo.GetByID(ctx, reservation.ID)
		if err != domain.ErrReservationNotFound {
			t.Errorf("got %v, want %v", err, domain.ErrReservationNotFound)
		}
	})
}

func TestDomainDateToPgDate(t *testing.T) {
	// This is used internally by the repository
	d := domain.Date{Year: 2026, Month: 5, Day: 15}
	pg := domainDateToPgDate(d)
	if !pg.Valid {
		t.Error("expected valid")
	}
	if pg.Time.Year() != 2026 || int(pg.Time.Month()) != 5 || pg.Time.Day() != 15 {
		t.Errorf("got %v, want 2026-05-15", pg.Time)
	}
}

func TestDateFromPgDate(t *testing.T) {
	pg := pgtype.Date{Time: time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC), Valid: true}
	d := dateFromPgDate(pg)
	if d != (domain.Date{Year: 2026, Month: 5, Day: 15}) {
		t.Errorf("got %+v, want {2026 5 15}", d)
	}
}

func TestReservationRepositoryInterface(t *testing.T) {
	var _ ports.ReservationRepository = (*ReservationRepository)(nil)
}
