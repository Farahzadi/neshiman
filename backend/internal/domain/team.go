package domain

import "github.com/google/uuid"

type Team struct {
	ID   uuid.UUID
	Name string
}

func NewTeam(name string) *Team {
	return &Team{
		ID:   uuid.New(),
		Name: name,
	}
}
