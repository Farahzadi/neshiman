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
		}, &mockUserRepo{})

		seat, err := svc.CreateSeat(context.Background(), roomID, teamID, "A1", 5, 3)
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

	t.Run("repo error", func(t *testing.T) {
		svc := NewSeatService(&mockSeatRepo{
			createFn: func(_ context.Context, _ *domain.Seat) error {
				return errors.New("db error")
			},
		}, &mockUserRepo{})
		_, err := svc.CreateSeat(context.Background(), roomID, teamID, "X", 1, 1)
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
	}, &mockUserRepo{})

	seat, err := svc.MoveSeat(context.Background(), id, 10, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if seat.Position.X != 10 || seat.Position.Y != 20 {
		t.Errorf("got %+v, want {10 20}", seat.Position)
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
	}, &mockUserRepo{})

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
	}, &mockUserRepo{})

	seats, err := svc.ListSeatsByRoom(context.Background(), roomID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(seats) != 2 {
		t.Errorf("got %d seats, want 2", len(seats))
	}
}

func TestSeatService_BulkSyncSeats(t *testing.T) {
	roomID := uuid.New()
	teamID := uuid.New()

	t.Run("success with mix of operations", func(t *testing.T) {
		existingID := uuid.New()
		svc := NewSeatService(&mockSeatRepo{
			bulkSyncFn: func(_ context.Context, gotID uuid.UUID, seats []domain.Seat) ([]domain.Seat, error) {
				if gotID != roomID {
					t.Errorf("got roomID %v, want %v", gotID, roomID)
				}
				if len(seats) != 2 {
					t.Errorf("got %d seats, want 2", len(seats))
				}
				return []domain.Seat{
					{ID: existingID, RoomID: roomID, TeamID: teamID, Label: "A1", Position: domain.Position{X: 0, Y: 0}},
					{ID: uuid.New(), RoomID: roomID, TeamID: teamID, Label: "B2", Position: domain.Position{X: 3, Y: 5}},
				}, nil
			},
		}, &mockUserRepo{})

		seats := []domain.Seat{
			{ID: existingID, TeamID: teamID, Label: "A1", Position: domain.Position{X: 0, Y: 0}},
			{TeamID: teamID, Label: "B2", Position: domain.Position{X: 3, Y: 5}},
		}
		result, err := svc.BulkSyncSeats(context.Background(), roomID, seats)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("got %d results, want 2", len(result))
		}
	})

	t.Run("repo error", func(t *testing.T) {
		svc := NewSeatService(&mockSeatRepo{
			bulkSyncFn: func(_ context.Context, _ uuid.UUID, _ []domain.Seat) ([]domain.Seat, error) {
				return nil, errors.New("db error")
			},
		}, &mockUserRepo{})
		seats := []domain.Seat{
			{TeamID: teamID, Label: "A1", Position: domain.Position{X: 0, Y: 0}},
		}
		_, err := svc.BulkSyncSeats(context.Background(), roomID, seats)
		if err == nil || err.Error() != "db error" {
			t.Errorf("got %v, want db error", err)
		}
	})
}

func TestSeatService_DeleteSeat(t *testing.T) {
	id := uuid.New()
	var deleted uuid.UUID
	svc := NewSeatService(&mockSeatRepo{
		deleteFn: func(_ context.Context, got uuid.UUID) error {
			deleted = got
			return nil
		},
	}, &mockUserRepo{})

	if err := svc.DeleteSeat(context.Background(), id); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deleted != id {
		t.Errorf("got %v, want %v", deleted, id)
	}
}

func TestSeatService_AssignUser(t *testing.T) {
	seatID := uuid.New()
	userID := uuid.New()
	superadminID := uuid.New()

	t.Run("success", func(t *testing.T) {
		svc := NewSeatService(&mockSeatRepo{
			getByAssignedUserFn: func(_ context.Context, _ uuid.UUID) (*domain.Seat, error) {
				return nil, domain.ErrSeatNotFound
			},
			getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Seat, error) {
				return &domain.Seat{ID: seatID}, nil
			},
			assignUserFn: func(_ context.Context, sid, uid uuid.UUID) (*domain.Seat, error) {
				if sid != seatID || uid != userID {
					t.Errorf("got seat=%v user=%v, want seat=%v user=%v", sid, uid, seatID, userID)
				}
				return &domain.Seat{ID: seatID, AssignedUserID: &uid}, nil
			},
		}, &mockUserRepo{
			getByIDFn: func(_ context.Context, id uuid.UUID) (*domain.User, error) {
				if id == superadminID {
					return &domain.User{ID: superadminID, Role: domain.RoleSuperAdmin}, nil
				}
				return &domain.User{ID: userID, Role: domain.RoleViewer}, nil
			},
		})

		seat, err := svc.AssignUser(context.Background(), superadminID, seatID, userID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if seat.AssignedUserID == nil || *seat.AssignedUserID != userID {
			t.Error("seat should be assigned to user")
		}
	})

	t.Run("user not found", func(t *testing.T) {
		svc := NewSeatService(&mockSeatRepo{
			getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Seat, error) {
				return &domain.Seat{ID: seatID}, nil
			},
		}, &mockUserRepo{
			getByIDFn: func(_ context.Context, id uuid.UUID) (*domain.User, error) {
				if id == superadminID {
					return &domain.User{ID: superadminID, Role: domain.RoleSuperAdmin}, nil
				}
				return nil, domain.ErrUserNotFound
			},
		})

		_, err := svc.AssignUser(context.Background(), superadminID, seatID, userID)
		if err != domain.ErrUserNotFound {
			t.Errorf("got %v, want %v", err, domain.ErrUserNotFound)
		}
	})

	t.Run("user already assigned to another seat", func(t *testing.T) {
		svc := NewSeatService(&mockSeatRepo{
			getByAssignedUserFn: func(_ context.Context, _ uuid.UUID) (*domain.Seat, error) {
				otherID := uuid.New()
				return &domain.Seat{ID: uuid.New(), AssignedUserID: &otherID}, nil
			},
		}, &mockUserRepo{
			getByIDFn: func(_ context.Context, id uuid.UUID) (*domain.User, error) {
				if id == superadminID {
					return &domain.User{ID: superadminID, Role: domain.RoleSuperAdmin}, nil
				}
				return &domain.User{ID: userID, Role: domain.RoleViewer}, nil
			},
		})

		_, err := svc.AssignUser(context.Background(), superadminID, seatID, userID)
		if err != domain.ErrUserAlreadyAssigned {
			t.Errorf("got %v, want %v", err, domain.ErrUserAlreadyAssigned)
		}
	})

	t.Run("seat not found", func(t *testing.T) {
		svc := NewSeatService(&mockSeatRepo{
			getByAssignedUserFn: func(_ context.Context, _ uuid.UUID) (*domain.Seat, error) {
				return nil, domain.ErrSeatNotFound
			},
			getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Seat, error) {
				return nil, domain.ErrSeatNotFound
			},
		}, &mockUserRepo{
			getByIDFn: func(_ context.Context, id uuid.UUID) (*domain.User, error) {
				if id == superadminID {
					return &domain.User{ID: superadminID, Role: domain.RoleSuperAdmin}, nil
				}
				return &domain.User{ID: userID, Role: domain.RoleViewer}, nil
			},
		})

		_, err := svc.AssignUser(context.Background(), superadminID, seatID, userID)
		if err != domain.ErrSeatNotFound {
			t.Errorf("got %v, want %v", err, domain.ErrSeatNotFound)
		}
	})
}

func TestSeatService_UnassignUser(t *testing.T) {
	seatID := uuid.New()
	superadminID := uuid.New()

	t.Run("success", func(t *testing.T) {
		svc := NewSeatService(&mockSeatRepo{
			unassignUserFn: func(_ context.Context, sid uuid.UUID) (*domain.Seat, error) {
				if sid != seatID {
					t.Errorf("got %v, want %v", sid, seatID)
				}
				return &domain.Seat{ID: seatID}, nil
			},
		}, &mockUserRepo{
			getByIDFn: func(_ context.Context, id uuid.UUID) (*domain.User, error) {
				return &domain.User{ID: id, Role: domain.RoleSuperAdmin}, nil
			},
		})

		seat, err := svc.UnassignUser(context.Background(), superadminID, seatID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if seat.AssignedUserID != nil {
			t.Error("seat should not have assigned user")
		}
	})

	t.Run("seat not found", func(t *testing.T) {
		svc := NewSeatService(&mockSeatRepo{
			unassignUserFn: func(_ context.Context, _ uuid.UUID) (*domain.Seat, error) {
				return nil, domain.ErrSeatNotFound
			},
		}, &mockUserRepo{
			getByIDFn: func(_ context.Context, id uuid.UUID) (*domain.User, error) {
				return &domain.User{ID: id, Role: domain.RoleSuperAdmin}, nil
			},
		})

		_, err := svc.UnassignUser(context.Background(), superadminID, seatID)
		if err != domain.ErrSeatNotFound {
			t.Errorf("got %v, want %v", err, domain.ErrSeatNotFound)
		}
	})
}
