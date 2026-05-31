package application

import (
	"context"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"

	"github.com/google/uuid"
)

type SeatService struct {
	seats ports.SeatRepository
	users ports.UserRepository
}

func NewSeatService(seats ports.SeatRepository, users ports.UserRepository) *SeatService {
	return &SeatService{seats: seats, users: users}
}

func (s *SeatService) CreateSeat(ctx context.Context, roomID, teamID uuid.UUID, label string, posX, posY int) (*domain.Seat, error) {
	pos := domain.Position{X: posX, Y: posY}
	seat, err := domain.NewSeat(roomID, teamID, label, pos)
	if err != nil {
		return nil, err
	}
	if err := s.seats.Create(ctx, seat); err != nil {
		return nil, err
	}
	return seat, nil
}

func (s *SeatService) GetSeat(ctx context.Context, id uuid.UUID) (*domain.Seat, error) {
	return s.seats.GetByID(ctx, id)
}

func (s *SeatService) ListSeatsByRoom(ctx context.Context, roomID uuid.UUID) ([]domain.Seat, error) {
	return s.seats.ListByRoom(ctx, roomID)
}

func (s *SeatService) MoveSeat(ctx context.Context, id uuid.UUID, posX, posY int) (*domain.Seat, error) {
	seat, err := s.seats.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	seat.MoveTo(domain.Position{X: posX, Y: posY})
	if err := s.seats.Update(ctx, seat); err != nil {
		return nil, err
	}
	return seat, nil
}

func (s *SeatService) DeleteSeat(ctx context.Context, id uuid.UUID) error {
	return s.seats.Delete(ctx, id)
}

func (s *SeatService) BulkSyncSeats(ctx context.Context, roomID uuid.UUID, seats []domain.Seat) ([]domain.Seat, error) {
	for i, seat := range seats {
		if _, err := domain.NewSeat(roomID, seat.TeamID, seat.Label, seat.Position); err != nil {
			return nil, err
		}
		seats[i].RoomID = roomID
	}
	return s.seats.BulkSync(ctx, roomID, seats)
}

// AssignUser permanently assigns a user to a seat (superadmin only).
// A user can only be assigned to one seat at a time.
func (s *SeatService) AssignUser(ctx context.Context, callerID, seatID, userID uuid.UUID) (*domain.Seat, error) {
	caller, err := s.users.GetByID(ctx, callerID)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	if !caller.IsSuperAdmin() {
		return nil, domain.ErrForbidden
	}

	targetUser, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Check user does not already have a permanent seat
	existing, err := s.seats.GetByAssignedUser(ctx, userID)
	if err != nil && err != domain.ErrSeatNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrUserAlreadyAssigned
	}

	seat, err := s.seats.GetByID(ctx, seatID)
	if err != nil {
		return nil, err
	}

	// Assign the seat
	seat, err = s.seats.AssignUser(ctx, seatID, targetUser.ID)
	if err != nil {
		return nil, err
	}
	return seat, nil
}

// UnassignUser removes the permanent user assignment from a seat (superadmin only).
func (s *SeatService) UnassignUser(ctx context.Context, callerID, seatID uuid.UUID) (*domain.Seat, error) {
	caller, err := s.users.GetByID(ctx, callerID)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	if !caller.IsSuperAdmin() {
		return nil, domain.ErrForbidden
	}

	seat, err := s.seats.UnassignUser(ctx, seatID)
	if err != nil {
		return nil, err
	}
	return seat, nil
}
