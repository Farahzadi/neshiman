package postgres

import (
	"context"
	"time"

	"neshiman/backend/db/sqlc"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationRepository struct {
	q    *sqlc.Queries
	pool *pgxpool.Pool
}

func NewReservationRepository(pool *pgxpool.Pool) ports.ReservationRepository {
	return &ReservationRepository{
		q:    sqlc.New(pool),
		pool: pool,
	}
}

func domainDateToTime(d domain.Date) time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
}

func (r *ReservationRepository) Create(ctx context.Context, reservation *domain.Reservation) error {
	result, err := r.q.CreateReservation(ctx, sqlc.CreateReservationParams{
		UserID: reservation.UserID,
		SeatID: reservation.SeatID,
		Date:   domainDateToTime(reservation.Date),
	})
	if err != nil {
		return err
	}
	reservation.ID = result.ID
	return nil
}

func (r *ReservationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Reservation, error) {
	result, err := r.q.GetReservationByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrReservationNotFound
		}
		return nil, err
	}
	return &domain.Reservation{
		ID:     result.ID,
		UserID: result.UserID,
		SeatID: result.SeatID,
		Date:   dateFromTime(result.Date),
	}, nil
}

func (r *ReservationRepository) GetBySeatAndDate(ctx context.Context, seatID uuid.UUID, date domain.Date) (*domain.Reservation, error) {
	result, err := r.q.GetReservationBySeatAndDate(ctx, sqlc.GetReservationBySeatAndDateParams{
		SeatID: seatID,
		Date:   domainDateToTime(date),
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &domain.Reservation{
		ID:     result.ID,
		UserID: result.UserID,
		SeatID: result.SeatID,
		Date:   dateFromTime(result.Date),
	}, nil
}

func (r *ReservationRepository) ListByUserAndWeek(ctx context.Context, userID uuid.UUID, start, end domain.Date) ([]domain.Reservation, error) {
	results, err := r.q.ListReservationsByUserAndWeek(ctx, sqlc.ListReservationsByUserAndWeekParams{
		UserID: userID,
		Date:   domainDateToTime(start),
		Date_2: domainDateToTime(end),
	})
	if err != nil {
		return nil, err
	}
	reservations := make([]domain.Reservation, len(results))
	for i, row := range results {
		reservations[i] = domain.Reservation{
			ID:     row.ID,
			UserID: row.UserID,
			SeatID: row.SeatID,
			Date:   dateFromTime(row.Date),
		}
	}
	return reservations, nil
}

func (r *ReservationRepository) CountByUserInWeek(ctx context.Context, userID uuid.UUID, start, end domain.Date) (int, error) {
	count, err := r.q.CountReservationsByUserInWeek(ctx, sqlc.CountReservationsByUserInWeekParams{
		UserID: userID,
		Date:   domainDateToTime(start),
		Date_2: domainDateToTime(end),
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *ReservationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteReservation(ctx, id)
}

func dateFromTime(t time.Time) domain.Date {
	return domain.Date{Year: t.Year(), Month: int(t.Month()), Day: t.Day()}
}
