package domain

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID        uuid.UUID
	Name      string
	DeletedAt *time.Time
}

func NewTeam(name string) *Team {
	return &Team{
		ID:   uuid.New(),
		Name: name,
	}
}
