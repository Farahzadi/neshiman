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
	daysSinceSat := (int(weekday) - int(time.Saturday) + 7) % 7
	start := t.AddDate(0, 0, -daysSinceSat)
	end := start.AddDate(0, 0, 6)
	return dateFromTime(start), dateFromTime(end)
}

func dateFromTime(t time.Time) domain.Date {
	return domain.Date{Year: t.Year(), Month: int(t.Month()), Day: t.Day()}
}

func (s *ReservationService) ReserveSeat(ctx context.Context, userID, seatID uuid.UUID, date domain.Date) (*domain.Reservation, error) {
	return s.reserveSeat(ctx, userID, seatID, date, false)
}

// AdminReserveSeat creates a reservation on behalf of another user.
// callerID is the admin, targetUserID is the person getting the seat.
func (s *ReservationService) AdminReserveSeat(ctx context.Context, callerID, targetUserID, seatID uuid.UUID, date domain.Date) (*domain.Reservation, error) {
	caller, err := s.users.GetByID(ctx, callerID)
	if err != nil {
		return nil, err
	}
	if !caller.IsSuperAdmin() && !caller.IsTeamAdmin() {
		return nil, domain.ErrForbidden
	}

	target, err := s.users.GetByID(ctx, targetUserID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// team_admin can only reserve for their own team
	if caller.IsTeamAdmin() {
		if target.TeamID == nil || caller.TeamID == nil || *target.TeamID != *caller.TeamID {
			return nil, domain.ErrWrongTeam
		}
	}

	return s.reserveSeat(ctx, targetUserID, seatID, date, true)
}

func (s *ReservationService) reserveSeat(ctx context.Context, userID, seatID uuid.UUID, date domain.Date, allowCrossTeam bool) (*domain.Reservation, error) {
	// Validate date is not in the past
	today := dateFromTime(time.Now())
	if date.Year < today.Year || (date.Year == today.Year && date.Month < today.Month) || (date.Year == today.Year && date.Month == today.Month && date.Day < today.Day) {
		return nil, domain.ErrPastDate
	}

	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	seat, err := s.seats.GetByID(ctx, seatID)
	if err != nil {
		return nil, domain.ErrSeatNotFound
	}

	// Check if seat is permanently assigned to another user
	if seat.IsPermanentlyAssigned() && *seat.AssignedUserID != userID {
		return nil, domain.ErrSeatPermanentlyAssigned
	}

	// Check if seat is already reserved for this date
	existing, _ := s.reservations.GetBySeatAndDate(ctx, seatID, date)
	if existing != nil {
		return nil, domain.ErrSeatAlreadyReserved
	}

	// Check cross-team: if seat is owned by a different team
	if user.TeamID != nil && seat.TeamID != *user.TeamID {
		if !allowCrossTeam {
			return nil, domain.ErrForbidden
		}
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

	// Self-cancel: anyone can cancel their own
	if reservation.UserID == userID {
		return s.reservations.Delete(ctx, reservationID)
	}

	// Admin cancel: check permissions
	caller, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if caller.IsSuperAdmin() {
		return s.reservations.Delete(ctx, reservationID)
	}
	if caller.IsTeamAdmin() {
		// Find the reservation's user to check team
		reservationUser, err := s.users.GetByID(ctx, reservation.UserID)
		if err != nil {
			return err
		}
		if reservationUser.TeamID != nil && caller.TeamID != nil && *reservationUser.TeamID == *caller.TeamID {
			return s.reservations.Delete(ctx, reservationID)
		}
	}

	return domain.ErrForbidden
}

func (s *ReservationService) ListByRoomAndDate(ctx context.Context, roomID uuid.UUID, date domain.Date) ([]domain.Reservation, error) {
	return s.reservations.ListByRoomAndDate(ctx, roomID, date)
}

func (s *ReservationService) ListByRoomAndDateWithDetails(ctx context.Context, roomID uuid.UUID, date domain.Date) ([]ports.ReservationWithDetails, error) {
	return s.reservations.ListByRoomAndDateWithDetails(ctx, roomID, date)
}

func (s *ReservationService) GetUserWeekReservations(ctx context.Context, userID uuid.UUID, date domain.Date) ([]domain.Reservation, error) {
	start, end := weekBounds(date)
	return s.reservations.ListByUserAndWeek(ctx, userID, start, end)
}

func (s *ReservationService) ListByDate(ctx context.Context, date domain.Date) ([]domain.Reservation, error) {
	return s.reservations.ListByDate(ctx, date)
}

func (s *ReservationService) ListByUserAndDate(ctx context.Context, userID uuid.UUID, date domain.Date) ([]domain.Reservation, error) {
	return s.reservations.ListByUserAndDate(ctx, userID, date)
}
