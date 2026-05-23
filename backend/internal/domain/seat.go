package domain

import (
	"time"

	"github.com/google/uuid"
)

type Seat struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	TeamID    uuid.UUID
	Label     string
	Position  Position
	Rotation  Rotation
	DeletedAt *time.Time
}

func NewSeat(roomID, teamID uuid.UUID, label string, pos Position, rotation Rotation) (*Seat, error) {
	if rotation != Rotation0 && rotation != Rotation90 && rotation != Rotation180 && rotation != Rotation270 {
		return nil, ErrInvalidRotation
	}
	return &Seat{
		ID:       uuid.New(),
		RoomID:   roomID,
		TeamID:   teamID,
		Label:    label,
		Position: pos,
		Rotation: rotation,
	}, nil
}

func (s *Seat) MoveTo(pos Position) {
	s.Position = pos
}

func (s *Seat) Rotate(rotation Rotation) {
	s.Rotation = rotation
}
