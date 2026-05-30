package postgres

import (
	"context"
	"neshiman/backend/db/sqlc"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	q    *sqlc.Queries
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) ports.UserRepository {
	return &UserRepository{
		q:    sqlc.New(pool),
		pool: pool,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	teamID := pgtype.UUID{Valid: false}
	if user.TeamID != nil {
		teamID = pgtype.UUID{Bytes: *user.TeamID, Valid: true}
	}
	email := pgtype.Text{Valid: false}
	if user.Email != "" {
		email = pgtype.Text{String: user.Email, Valid: true}
	}
	result, err := r.q.CreateUser(ctx, sqlc.CreateUserParams{
		Name:         user.Name,
		Email:        email,
		TeamID:       teamID,
		Role:         string(user.Role),
		WeeklyLimit:  int32(user.WeeklyLimit),
		PasswordHash: user.PasswordHash,
	})
	if err != nil {
		return err
	}
	user.ID = result.ID
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	result, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return mapUser(result), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	result, err := r.q.GetUserByEmail(ctx, pgtype.Text{String: email, Valid: email != ""})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return mapUser(result), nil
}

func (r *UserRepository) GetByName(ctx context.Context, name string) (*domain.User, error) {
	result, err := r.q.GetUserByName(ctx, name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return mapUser(result), nil
}

func (r *UserRepository) ListByTeam(ctx context.Context, teamID uuid.UUID) ([]domain.User, error) {
	results, err := r.q.ListUsersByTeam(ctx, pgtype.UUID{Bytes: teamID, Valid: true})
	if err != nil {
		return nil, err
	}
	users := make([]domain.User, len(results))
	for i, row := range results {
		users[i] = *mapUser(row)
	}
	return users, nil
}

func (r *UserRepository) ListAll(ctx context.Context) ([]domain.User, error) {
	results, err := r.q.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]domain.User, len(results))
	for i, row := range results {
		users[i] = *mapUser(row)
	}
	return users, nil
}

func (r *UserRepository) UpdateWeeklyLimit(ctx context.Context, userID uuid.UUID, limit domain.WeeklyLimit) error {
	_, err := r.q.UpdateUserWeeklyLimit(ctx, sqlc.UpdateUserWeeklyLimitParams{
		ID:          userID,
		WeeklyLimit: int32(limit),
	})
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	return r.q.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
		ID:           userID,
		PasswordHash: passwordHash,
	})
}

func (r *UserRepository) UpdateRole(ctx context.Context, userID uuid.UUID, role domain.Role) error {
	return r.q.UpdateUserRole(ctx, sqlc.UpdateUserRoleParams{
		ID:   userID,
		Role: string(role),
	})
}

func (r *UserRepository) UpdateTeamID(ctx context.Context, userID uuid.UUID, teamID *uuid.UUID) error {
	t := pgtype.UUID{Valid: false}
	if teamID != nil {
		t = pgtype.UUID{Bytes: *teamID, Valid: true}
	}
	return r.q.UpdateUserTeamID(ctx, sqlc.UpdateUserTeamIDParams{
		ID:     userID,
		TeamID: t,
	})
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.SoftDeleteUser(ctx, id)
}

func mapUser(row sqlc.User) *domain.User {
	user := &domain.User{
		ID:           row.ID,
		Name:         row.Name,
		PasswordHash: row.PasswordHash,
		Role:         domain.Role(row.Role),
		WeeklyLimit:  domain.WeeklyLimit(row.WeeklyLimit),
	}
	if row.Email.Valid {
		user.Email = row.Email.String
	}
	if row.TeamID.Valid {
		t := uuid.UUID(row.TeamID.Bytes)
		user.TeamID = &t
	}
	return user
}
