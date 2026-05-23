package application

import (
	"context"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"

	"github.com/google/uuid"
)

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
