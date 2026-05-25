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
	DeletedAt *time.Time
}

func NewSeat(roomID, teamID uuid.UUID, label string, pos Position) (*Seat, error) {
	return &Seat{
		ID:       uuid.New(),
		RoomID:   roomID,
		TeamID:   teamID,
		Label:    label,
		Position: pos,
	}, nil
}

func (s *Seat) MoveTo(pos Position) {
	s.Position = pos
}
