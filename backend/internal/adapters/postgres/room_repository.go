package postgres

import (
	"context"
	"neshiman/backend/db/sqlc"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomRepository struct {
	q    *sqlc.Queries
	pool *pgxpool.Pool
}

func NewRoomRepository(pool *pgxpool.Pool) ports.RoomRepository {
	return &RoomRepository{
		q:    sqlc.New(pool),
		pool: pool,
	}
}

func (r *RoomRepository) Create(ctx context.Context, room *domain.Room) error {
	result, err := r.q.CreateRoom(ctx, sqlc.CreateRoomParams{
		Name:       room.Name,
		GridWidth:  int32(room.GridWidth),
		GridHeight: int32(room.GridHeight),
	})
	if err != nil {
		return err
	}
	room.ID = result.ID
	room.Name = result.Name
	room.GridWidth = int(result.GridWidth)
	room.GridHeight = int(result.GridHeight)
	return nil
}

func (r *RoomRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	result, err := r.q.GetRoomByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrRoomNotFound
		}
		return nil, err
	}
	return &domain.Room{
		ID:         result.ID,
		Name:       result.Name,
		GridWidth:  int(result.GridWidth),
		GridHeight: int(result.GridHeight),
	}, nil
}

func (r *RoomRepository) List(ctx context.Context) ([]domain.Room, error) {
	results, err := r.q.ListRooms(ctx)
	if err != nil {
		return nil, err
	}
	rooms := make([]domain.Room, len(results))
	for i, row := range results {
		rooms[i] = domain.Room{
			ID:         row.ID,
			Name:       row.Name,
			GridWidth:  int(row.GridWidth),
			GridHeight: int(row.GridHeight),
		}
	}
	return rooms, nil
}

func (r *RoomRepository) Update(ctx context.Context, room *domain.Room) error {
	_, err := r.q.UpdateRoom(ctx, sqlc.UpdateRoomParams{
		ID:         room.ID,
		Name:       room.Name,
		GridWidth:  int32(room.GridWidth),
		GridHeight: int32(room.GridHeight),
	})
	return err
}

func (r *RoomRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.SoftDeleteRoom(ctx, id)
}
