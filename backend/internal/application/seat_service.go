package application

import (
	"context"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"

	"github.com/google/uuid"
)

type SeatService struct {
	seats ports.SeatRepository
}

func NewSeatService(seats ports.SeatRepository) *SeatService {
	return &SeatService{seats: seats}
}

func (s *SeatService) CreateSeat(ctx context.Context, roomID, teamID uuid.UUID, label string, posX, posY int, rotation int) (*domain.Seat, error) {
	pos := domain.Position{X: posX, Y: posY}
	rot := domain.Rotation(rotation)
	seat, err := domain.NewSeat(roomID, teamID, label, pos, rot)
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

func (s *SeatService) RotateSeat(ctx context.Context, id uuid.UUID, rotation int) (*domain.Seat, error) {
	seat, err := s.seats.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	rot := domain.Rotation(rotation)
	if rot != domain.Rotation0 && rot != domain.Rotation90 && rot != domain.Rotation180 && rot != domain.Rotation270 {
		return nil, domain.ErrInvalidRotation
	}
	seat.Rotate(rot)
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
		if _, err := domain.NewSeat(roomID, seat.TeamID, seat.Label, seat.Position, seat.Rotation); err != nil {
			return nil, err
		}
		seats[i].RoomID = roomID
	}
	return s.seats.BulkSync(ctx, roomID, seats)
}
