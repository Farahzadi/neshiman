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
	result, err := r.q.CreateUser(ctx, sqlc.CreateUserParams{
		Name:        user.Name,
		Email:       user.Email,
		TeamID:      teamID,
		Role:        string(user.Role),
		WeeklyLimit: int32(user.WeeklyLimit),
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
	user := &domain.User{
		ID:          result.ID,
		Name:        result.Name,
		Email:       result.Email,
		Role:        domain.Role(result.Role),
		WeeklyLimit: domain.WeeklyLimit(result.WeeklyLimit),
	}
	if result.TeamID.Valid {
		t := uuid.UUID(result.TeamID.Bytes)
		user.TeamID = &t
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	result, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	user := &domain.User{
		ID:          result.ID,
		Name:        result.Name,
		Email:       result.Email,
		Role:        domain.Role(result.Role),
		WeeklyLimit: domain.WeeklyLimit(result.WeeklyLimit),
	}
	if result.TeamID.Valid {
		t := uuid.UUID(result.TeamID.Bytes)
		user.TeamID = &t
	}
	return user, nil
}

func (r *UserRepository) ListByTeam(ctx context.Context, teamID uuid.UUID) ([]domain.User, error) {
	results, err := r.q.ListUsersByTeam(ctx, pgtype.UUID{Bytes: teamID, Valid: true})
	if err != nil {
		return nil, err
	}
	users := make([]domain.User, len(results))
	for i, row := range results {
		users[i] = domain.User{
			ID:          row.ID,
			Name:        row.Name,
			Email:       row.Email,
			Role:        domain.Role(row.Role),
			WeeklyLimit: domain.WeeklyLimit(row.WeeklyLimit),
		}
		if row.TeamID.Valid {
			t := uuid.UUID(row.TeamID.Bytes)
			users[i].TeamID = &t
		}
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
		users[i] = domain.User{
			ID:          row.ID,
			Name:        row.Name,
			Email:       row.Email,
			Role:        domain.Role(row.Role),
			WeeklyLimit: domain.WeeklyLimit(row.WeeklyLimit),
		}
		if row.TeamID.Valid {
			t := uuid.UUID(row.TeamID.Bytes)
			users[i].TeamID = &t
		}
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

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteUser(ctx, id)
}
