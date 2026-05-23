package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	TeamID       *uuid.UUID
	Role         Role
	WeeklyLimit  WeeklyLimit
	DeletedAt    *time.Time
}

func NewUser(name, email string, teamID *uuid.UUID, role Role, weeklyLimit WeeklyLimit) *User {
	if weeklyLimit == 0 {
		weeklyLimit = 2
	}
	return &User{
		ID:          uuid.New(),
		Name:        name,
		Email:       email,
		TeamID:      teamID,
		Role:        role,
		WeeklyLimit: weeklyLimit,
	}
}

func (u *User) IsSuperAdmin() bool {
	return u.Role == RoleSuperAdmin
}

func (u *User) IsTeamAdmin() bool {
	return u.Role == RoleTeamAdmin
}

func (u *User) IsAdminOfTeam(teamID uuid.UUID) bool {
	return u.IsSuperAdmin() || (u.IsTeamAdmin() && u.TeamID != nil && *u.TeamID == teamID)
}
