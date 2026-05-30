package application

import (
	"context"
	"errors"
	"neshiman/backend/internal/domain"
	"testing"

	"github.com/google/uuid"
)

func TestReservationService_ReserveSeat(t *testing.T) {
	userID := uuid.New()
	seatID := uuid.New()
	teamID := uuid.New()
	date := domain.Date{Year: 2027, Month: 6, Day: 15}

	t.Run("success", func(t *testing.T) {
		var created *domain.Reservation
		svc := NewReservationService(
			&mockReservationRepo{
				getBySeatAndDateFn: func(_ context.Context, _ uuid.UUID, _ domain.Date) (*domain.Reservation, error) {
					return nil, nil
				},
				countByUserInWeekFn: func(_ context.Context, _ uuid.UUID, _, _ domain.Date) (int, error) {
					return 0, nil
				},
				createFn: func(_ context.Context, r *domain.Reservation) error {
					created = r
					return nil
				},
			},
			&mockSeatRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Seat, error) {
					return &domain.Seat{ID: seatID, RoomID: uuid.New(), TeamID: teamID}, nil
				},
			},
			&mockUserRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.User, error) {
					return &domain.User{ID: userID, TeamID: &teamID, WeeklyLimit: 2}, nil
				},
			},
			&mockTxManager{
				withinTxFn: func(_ context.Context, fn func(context.Context) error) error {
					return fn(context.Background())
				},
			},
		)

		r, err := svc.ReserveSeat(context.Background(), userID, seatID, date)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if r.UserID != userID || r.SeatID != seatID || r.Date != date {
			t.Error("reservation fields mismatch")
		}
		if created != r {
			t.Error("created reservation should be same pointer")
		}
	})

	t.Run("user not found", func(t *testing.T) {
		svc := NewReservationService(
			&mockReservationRepo{},
			&mockSeatRepo{},
			&mockUserRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.User, error) {
					return nil, domain.ErrUserNotFound
				},
			},
			&mockTxManager{},
		)
		_, err := svc.ReserveSeat(context.Background(), userID, seatID, date)
		if !errors.Is(err, domain.ErrUserNotFound) {
			t.Errorf("got %v, want %v", err, domain.ErrUserNotFound)
		}
	})

	t.Run("seat not found", func(t *testing.T) {
		svc := NewReservationService(
			&mockReservationRepo{},
			&mockSeatRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Seat, error) {
					return nil, domain.ErrSeatNotFound
				},
			},
			&mockUserRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.User, error) {
					return &domain.User{ID: userID, TeamID: &teamID, WeeklyLimit: 2}, nil
				},
			},
			&mockTxManager{},
		)
		_, err := svc.ReserveSeat(context.Background(), userID, seatID, date)
		if !errors.Is(err, domain.ErrSeatNotFound) {
			t.Errorf("got %v, want %v", err, domain.ErrSeatNotFound)
		}
	})

	t.Run("seat already reserved", func(t *testing.T) {
		svc := NewReservationService(
			&mockReservationRepo{
				getBySeatAndDateFn: func(_ context.Context, _ uuid.UUID, _ domain.Date) (*domain.Reservation, error) {
					return &domain.Reservation{}, nil
				},
			},
			&mockSeatRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Seat, error) {
					return &domain.Seat{ID: seatID, TeamID: teamID}, nil
				},
			},
			&mockUserRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.User, error) {
					return &domain.User{ID: userID, TeamID: &teamID, WeeklyLimit: 2}, nil
				},
			},
			&mockTxManager{},
		)
		_, err := svc.ReserveSeat(context.Background(), userID, seatID, date)
		if !errors.Is(err, domain.ErrSeatAlreadyReserved) {
			t.Errorf("got %v, want %v", err, domain.ErrSeatAlreadyReserved)
		}
	})

	t.Run("cross-team forbidden", func(t *testing.T) {
		otherTeamID := uuid.New()
		svc := NewReservationService(
			&mockReservationRepo{
				getBySeatAndDateFn: func(_ context.Context, _ uuid.UUID, _ domain.Date) (*domain.Reservation, error) {
					return nil, nil
				},
			},
			&mockSeatRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Seat, error) {
					return &domain.Seat{ID: seatID, TeamID: otherTeamID}, nil
				},
			},
			&mockUserRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.User, error) {
					return &domain.User{ID: userID, TeamID: &teamID, WeeklyLimit: 2}, nil
				},
			},
			&mockTxManager{},
		)
		_, err := svc.ReserveSeat(context.Background(), userID, seatID, date)
		if !errors.Is(err, domain.ErrForbidden) {
			t.Errorf("got %v, want %v", err, domain.ErrForbidden)
		}
	})

	t.Run("no team restriction", func(t *testing.T) {
		svc := NewReservationService(
			&mockReservationRepo{
				getBySeatAndDateFn: func(_ context.Context, _ uuid.UUID, _ domain.Date) (*domain.Reservation, error) {
					return nil, nil
				},
				countByUserInWeekFn: func(_ context.Context, _ uuid.UUID, _, _ domain.Date) (int, error) {
					return 0, nil
				},
				createFn: func(_ context.Context, _ *domain.Reservation) error {
					return nil
				},
			},
			&mockSeatRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Seat, error) {
					return &domain.Seat{ID: seatID, TeamID: uuid.New()}, nil
				},
			},
			&mockUserRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.User, error) {
					return &domain.User{ID: userID, TeamID: nil, WeeklyLimit: 2}, nil
				},
			},
			&mockTxManager{
				withinTxFn: func(_ context.Context, fn func(context.Context) error) error {
					return fn(context.Background())
				},
			},
		)
		_, err := svc.ReserveSeat(context.Background(), userID, seatID, date)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("weekly limit exceeded", func(t *testing.T) {
		svc := NewReservationService(
			&mockReservationRepo{
				getBySeatAndDateFn: func(_ context.Context, _ uuid.UUID, _ domain.Date) (*domain.Reservation, error) {
					return nil, nil
				},
				countByUserInWeekFn: func(_ context.Context, _ uuid.UUID, _, _ domain.Date) (int, error) {
					return 2, nil
				},
			},
			&mockSeatRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Seat, error) {
					return &domain.Seat{ID: seatID, TeamID: teamID}, nil
				},
			},
			&mockUserRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.User, error) {
					return &domain.User{ID: userID, TeamID: &teamID, WeeklyLimit: 2}, nil
				},
			},
			&mockTxManager{},
		)
		_, err := svc.ReserveSeat(context.Background(), userID, seatID, date)
		if !errors.Is(err, domain.ErrWeeklyLimitExceeded) {
			t.Errorf("got %v, want %v", err, domain.ErrWeeklyLimitExceeded)
		}
	})

	t.Run("tx error", func(t *testing.T) {
		svc := NewReservationService(
			&mockReservationRepo{
				getBySeatAndDateFn: func(_ context.Context, _ uuid.UUID, _ domain.Date) (*domain.Reservation, error) {
					return nil, nil
				},
				countByUserInWeekFn: func(_ context.Context, _ uuid.UUID, _, _ domain.Date) (int, error) {
					return 0, nil
				},
			},
			&mockSeatRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Seat, error) {
					return &domain.Seat{ID: seatID, TeamID: teamID}, nil
				},
			},
			&mockUserRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.User, error) {
					return &domain.User{ID: userID, TeamID: &teamID, WeeklyLimit: 2}, nil
				},
			},
			&mockTxManager{
				withinTxFn: func(_ context.Context, _ func(context.Context) error) error {
					return errors.New("tx error")
				},
			},
		)
		_, err := svc.ReserveSeat(context.Background(), userID, seatID, date)
		if err == nil || err.Error() != "tx error" {
			t.Errorf("got %v, want tx error", err)
		}
	})
}

func TestReservationService_CancelReservation(t *testing.T) {
	userID := uuid.New()
	reservationID := uuid.New()

	t.Run("success", func(t *testing.T) {
		var deleted uuid.UUID
		svc := NewReservationService(
			&mockReservationRepo{
				getByIDFn: func(_ context.Context, id uuid.UUID) (*domain.Reservation, error) {
					return &domain.Reservation{ID: id, UserID: userID}, nil
				},
				deleteFn: func(_ context.Context, id uuid.UUID) error {
					deleted = id
					return nil
				},
			},
			&mockSeatRepo{},
			&mockUserRepo{},
			&mockTxManager{},
		)
		err := svc.CancelReservation(context.Background(), reservationID, userID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if deleted != reservationID {
			t.Errorf("got %v, want %v", deleted, reservationID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := NewReservationService(
			&mockReservationRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Reservation, error) {
					return nil, domain.ErrReservationNotFound
				},
			},
			&mockSeatRepo{},
			&mockUserRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.User, error) {
					return &domain.User{Role: domain.RoleViewer}, nil
				},
			},
			&mockTxManager{},
		)
		err := svc.CancelReservation(context.Background(), reservationID, userID)
		if !errors.Is(err, domain.ErrReservationNotFound) {
			t.Errorf("got %v, want %v", err, domain.ErrReservationNotFound)
		}
	})

	t.Run("wrong user", func(t *testing.T) {
		callerUserID := uuid.New()
		svc := NewReservationService(
			&mockReservationRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Reservation, error) {
					return &domain.Reservation{ID: reservationID, UserID: uuid.New()}, nil
				},
			},
			&mockSeatRepo{},
			&mockUserRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.User, error) {
					return &domain.User{ID: callerUserID, Role: domain.RoleViewer, TeamID: &uuid.UUID{}}, nil
				},
			},
			&mockTxManager{},
		)
		err := svc.CancelReservation(context.Background(), reservationID, callerUserID)
		if !errors.Is(err, domain.ErrForbidden) {
			t.Errorf("got %v, want %v", err, domain.ErrForbidden)
		}
	})
}

func TestReservationService_GetUserWeekReservations(t *testing.T) {
	userID := uuid.New()
	date := domain.Date{Year: 2026, Month: 5, Day: 20}

	svc := NewReservationService(
		&mockReservationRepo{
			listByUserAndWeekFn: func(_ context.Context, got uuid.UUID, start, end domain.Date) ([]domain.Reservation, error) {
				if got != userID {
					t.Errorf("got user %v, want %v", got, userID)
				}
				return []domain.Reservation{{UserID: userID}}, nil
			},
		},
		&mockSeatRepo{},
		&mockUserRepo{},
		&mockTxManager{},
	)

	rs, err := svc.GetUserWeekReservations(context.Background(), userID, date)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rs) != 1 {
		t.Errorf("got %d, want 1", len(rs))
	}
}

func TestWeekBounds(t *testing.T) {
	tests := []struct {
		name         string
		input        domain.Date
		wantStart    domain.Date
		wantEnd      domain.Date
	}{
		{
			name:      "monday",
			input:     domain.Date{Year: 2026, Month: 5, Day: 18},
			wantStart: domain.Date{Year: 2026, Month: 5, Day: 18},
			wantEnd:   domain.Date{Year: 2026, Month: 5, Day: 24},
		},
		{
			name:      "wednesday",
			input:     domain.Date{Year: 2026, Month: 5, Day: 20},
			wantStart: domain.Date{Year: 2026, Month: 5, Day: 18},
			wantEnd:   domain.Date{Year: 2026, Month: 5, Day: 24},
		},
		{
			name:      "sunday",
			input:     domain.Date{Year: 2026, Month: 5, Day: 24},
			wantStart: domain.Date{Year: 2026, Month: 5, Day: 18},
			wantEnd:   domain.Date{Year: 2026, Month: 5, Day: 24},
		},
		{
			name:      "cross month boundary",
			input:     domain.Date{Year: 2026, Month: 5, Day: 29},
			wantStart: domain.Date{Year: 2026, Month: 5, Day: 25},
			wantEnd:   domain.Date{Year: 2026, Month: 5, Day: 31},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd := weekBounds(tt.input)
			if gotStart != tt.wantStart {
				t.Errorf("start: got %+v, want %+v", gotStart, tt.wantStart)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("end: got %+v, want %+v", gotEnd, tt.wantEnd)
			}
		})
	}
}
