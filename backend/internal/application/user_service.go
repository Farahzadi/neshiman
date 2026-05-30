package application

import (
	"context"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const DefaultPassword = "password"

type UserService struct {
	users ports.UserRepository
}

func NewUserService(users ports.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) CreateUser(ctx context.Context, name, email string, teamID *uuid.UUID, role domain.Role, weeklyLimit domain.WeeklyLimit) (*domain.User, error) {
	if weeklyLimit > domain.MaxWeeklyLimit {
		return nil, domain.ErrWeeklyLimitTooHigh
	}
	user := domain.NewUser(name, email, teamID, role, weeklyLimit)
	hash, err := bcrypt.GenerateFromPassword([]byte(DefaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = string(hash)
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.users.GetByEmail(ctx, email)
}

func (s *UserService) ListUsersByTeam(ctx context.Context, teamID uuid.UUID) ([]domain.User, error) {
	return s.users.ListByTeam(ctx, teamID)
}

func (s *UserService) ListAllUsers(ctx context.Context) ([]domain.User, error) {
	return s.users.ListAll(ctx)
}

func (s *UserService) UpdateWeeklyLimit(ctx context.Context, userID uuid.UUID, limit domain.WeeklyLimit) error {
	if limit > domain.MaxWeeklyLimit {
		return domain.ErrWeeklyLimitTooHigh
	}
	return s.users.UpdateWeeklyLimit(ctx, userID, limit)
}

func (s *UserService) UpdateUserRole(ctx context.Context, callerID, targetID uuid.UUID, role domain.Role) error {
	if role != domain.RoleSuperAdmin && role != domain.RoleTeamAdmin && role != domain.RoleViewer {
		return domain.ErrForbidden
	}

	caller, err := s.users.GetByID(ctx, callerID)
	if err != nil {
		return err
	}

	target, err := s.users.GetByID(ctx, targetID)
	if err != nil {
		return err
	}

	if target.IsSuperAdmin() {
		return domain.ErrCannotChangeSuperAdmin
	}

	// team_admin can only change viewers in their team
	if caller.IsTeamAdmin() {
		if !caller.IsAdminOfTeam(*target.TeamID) {
			return domain.ErrWrongTeam
		}
		if target.Role != domain.RoleViewer {
			return domain.ErrForbidden
		}
	}

	return s.users.UpdateRole(ctx, targetID, role)
}

func (s *UserService) UpdateUserTeam(ctx context.Context, callerID, targetID uuid.UUID, teamID *uuid.UUID) error {
	caller, err := s.users.GetByID(ctx, callerID)
	if err != nil {
		return err
	}
	if !caller.IsSuperAdmin() {
		return domain.ErrForbidden
	}

	target, err := s.users.GetByID(ctx, targetID)
	if err != nil {
		return err
	}
	if target.IsSuperAdmin() {
		return domain.ErrCannotChangeSuperAdmin
	}

	// If teamID is provided, verify team exists
	if teamID != nil {
		// We'll validate team existence via the handler
	}

	return s.users.UpdateTeamID(ctx, targetID, teamID)
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user.IsSuperAdmin() {
		return domain.ErrCannotDeleteSuperAdmin
	}
	return s.users.Delete(ctx, id)
}
