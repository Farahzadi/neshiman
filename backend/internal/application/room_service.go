package application

import (
	"context"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"

	"github.com/google/uuid"
)

type RoomService struct {
	rooms ports.RoomRepository
}

func NewRoomService(rooms ports.RoomRepository) *RoomService {
	return &RoomService{rooms: rooms}
}

func (s *RoomService) CreateRoom(ctx context.Context, name string, gridWidth, gridHeight int) (*domain.Room, error) {
	room, err := domain.NewRoom(name, gridWidth, gridHeight)
	if err != nil {
		return nil, err
	}
	if err := s.rooms.Create(ctx, room); err != nil {
		return nil, err
	}
	return room, nil
}

func (s *RoomService) GetRoom(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	return s.rooms.GetByID(ctx, id)
}

func (s *RoomService) ListRooms(ctx context.Context) ([]domain.Room, error) {
	return s.rooms.List(ctx)
}

func (s *RoomService) UpdateRoom(ctx context.Context, id uuid.UUID, name string, gridWidth, gridHeight int) (*domain.Room, error) {
	room, err := s.rooms.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	room.Name = name
	room.GridWidth = gridWidth
	room.GridHeight = gridHeight
	if err := s.rooms.Update(ctx, room); err != nil {
		return nil, err
	}
	return room, nil
}

func (s *RoomService) DeleteRoom(ctx context.Context, id uuid.UUID) error {
	return s.rooms.Delete(ctx, id)
}
