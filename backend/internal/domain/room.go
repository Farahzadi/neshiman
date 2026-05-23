package domain

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID         uuid.UUID
	Name       string
	GridWidth  int
	GridHeight int
	DeletedAt  *time.Time
}

func NewRoom(name string, gridWidth, gridHeight int) (*Room, error) {
	if gridWidth <= 0 || gridHeight <= 0 {
		return nil, ErrInvalidGridSize
	}
	return &Room{
		ID:         uuid.New(),
		Name:       name,
		GridWidth:  gridWidth,
		GridHeight: gridHeight,
	}, nil
}
