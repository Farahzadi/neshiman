package application

import (
	"context"
	"time"

	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users     ports.UserRepository
	teams     ports.TeamRepository
	jwtSecret string
}

func NewAuthService(users ports.UserRepository, teams ports.TeamRepository, jwtSecret string) *AuthService {
	return &AuthService{users: users, teams: teams, jwtSecret: jwtSecret}
}

type LoginResult struct {
	Token string
	User  *domain.User
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	user, err := s.users.GetByName(ctx, username)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  user.ID.String(),
		"role": string(user.Role),
		"iat":  now.Unix(),
		"exp":  now.Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &LoginResult{Token: signed, User: user}, nil
}

func (s *AuthService) SetPassword(ctx context.Context, callerID, targetID uuid.UUID, password string) error {
	if password == "" {
		return domain.ErrPasswordRequired
	}

	// Self-service: anyone can set their own password
	if callerID == targetID {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		return s.users.UpdatePassword(ctx, targetID, string(hash))
	}

	// Admin override: team_admin+ can set passwords
	caller, err := s.users.GetByID(ctx, callerID)
	if err != nil {
		return err
	}
	if !caller.IsSuperAdmin() && !caller.IsTeamAdmin() {
		return domain.ErrForbidden
	}

	target, err := s.users.GetByID(ctx, targetID)
	if err != nil {
		return err
	}

	// team_admin can only set passwords for their own team members
	if caller.IsTeamAdmin() {
		if target.TeamID == nil || caller.TeamID == nil || *target.TeamID != *caller.TeamID {
			return domain.ErrWrongTeam
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.users.UpdatePassword(ctx, targetID, string(hash))
}
