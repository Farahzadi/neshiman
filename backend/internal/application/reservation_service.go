package application

import (
	"context"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"
	"time"

	"github.com/google/uuid"
)

type ReservationService struct {
	reservations ports.ReservationRepository
	seats        ports.SeatRepository
	users        ports.UserRepository
	tx           ports.TxManager
}

func NewReservationService(
	reservations ports.ReservationRepository,
	seats ports.SeatRepository,
	users ports.UserRepository,
	tx ports.TxManager,
) *ReservationService {
	return &ReservationService{
		reservations: reservations,
		seats:        seats,
		users:        users,
		tx:           tx,
	}
}

func weekBounds(date domain.Date) (domain.Date, domain.Date) {
	t := time.Date(date.Year, time.Month(date.Month), date.Day, 0, 0, 0, 0, time.UTC)
	weekday := t.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	start := t.AddDate(0, 0, -int(weekday-time.Monday))
	end := start.AddDate(0, 0, 6)
	return dateFromTime(start), dateFromTime(end)
}

func dateFromTime(t time.Time) domain.Date {
	return domain.Date{Year: t.Year(), Month: int(t.Month()), Day: t.Day()}
}

func (s *ReservationService) ReserveSeat(ctx context.Context, userID, seatID uuid.UUID, date domain.Date) (*domain.Reservation, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	seat, err := s.seats.GetByID(ctx, seatID)
	if err != nil {
		return nil, domain.ErrSeatNotFound
	}

	// Check if seat is already reserved for this date
	existing, _ := s.reservations.GetBySeatAndDate(ctx, seatID, date)
	if existing != nil {
		return nil, domain.ErrSeatAlreadyReserved
	}

	// Check cross-team: if seat is owned by a different team, verify approval
	if user.TeamID != nil && seat.TeamID != *user.TeamID {
		// cross-team reservation requires approved request
		return nil, domain.ErrForbidden
	}

	// Check weekly limit
	start, end := weekBounds(date)
	count, err := s.reservations.CountByUserInWeek(ctx, userID, start, end)
	if err != nil {
		return nil, err
	}
	if count >= int(user.WeeklyLimit) {
		return nil, domain.ErrWeeklyLimitExceeded
	}

	var reservation *domain.Reservation
	err = s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		reservation = domain.NewReservation(userID, seatID, date)
		return s.reservations.Create(txCtx, reservation)
	})
	if err != nil {
		return nil, err
	}
	return reservation, nil
}

func (s *ReservationService) CancelReservation(ctx context.Context, reservationID, userID uuid.UUID) error {
	reservation, err := s.reservations.GetByID(ctx, reservationID)
	if err != nil {
		return domain.ErrReservationNotFound
	}
	if reservation.UserID != userID {
		return domain.ErrForbidden
	}
	return s.reservations.Delete(ctx, reservationID)
}

func (s *ReservationService) GetUserWeekReservations(ctx context.Context, userID uuid.UUID, date domain.Date) ([]domain.Reservation, error) {
	start, end := weekBounds(date)
	return s.reservations.ListByUserAndWeek(ctx, userID, start, end)
}
