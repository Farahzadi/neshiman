package application

import (
	"context"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"

	"github.com/google/uuid"
)

type CrossTeamRequestService struct {
	requests    ports.CrossTeamRequestRepository
	users       ports.UserRepository
	seats       ports.SeatRepository
	reservations ports.ReservationRepository
	tx          ports.TxManager
}

func NewCrossTeamRequestService(
	requests ports.CrossTeamRequestRepository,
	users ports.UserRepository,
	seats ports.SeatRepository,
	reservations ports.ReservationRepository,
	tx ports.TxManager,
) *CrossTeamRequestService {
	return &CrossTeamRequestService{
		requests:    requests,
		users:       users,
		seats:       seats,
		reservations: reservations,
		tx:          tx,
	}
}

func (s *CrossTeamRequestService) CreateRequest(ctx context.Context, requestingUserID, targetSeatID uuid.UUID, date domain.Date) (*domain.CrossTeamRequest, error) {
	req := domain.NewCrossTeamRequest(requestingUserID, targetSeatID, date)
	if err := s.requests.Create(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *CrossTeamRequestService) GetRequest(ctx context.Context, id uuid.UUID) (*domain.CrossTeamRequest, error) {
	return s.requests.GetByID(ctx, id)
}

func (s *CrossTeamRequestService) ListByStatus(ctx context.Context, status domain.RequestStatus) ([]domain.CrossTeamRequest, error) {
	return s.requests.ListByStatus(ctx, status)
}

func (s *CrossTeamRequestService) ListPendingByTeam(ctx context.Context, teamID uuid.UUID) ([]domain.CrossTeamRequest, error) {
	return s.requests.ListPendingByTeam(ctx, teamID)
}

func (s *CrossTeamRequestService) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.CrossTeamRequest, error) {
	return s.requests.ListByUser(ctx, userID)
}

func (s *CrossTeamRequestService) ListAll(ctx context.Context) ([]domain.CrossTeamRequest, error) {
	return s.requests.ListAll(ctx)
}

// canModifyRequest checks if the caller is authorized to approve/reject this request.
// The caller must be a team_admin of the seat's owning team, or a superadmin.
func (s *CrossTeamRequestService) canModifyRequest(ctx context.Context, callerID uuid.UUID, req *domain.CrossTeamRequest) error {
	caller, err := s.users.GetByID(ctx, callerID)
	if err != nil {
		return err
	}

	if caller.IsSuperAdmin() {
		return nil
	}

	seat, err := s.seats.GetByID(ctx, req.TargetSeatID)
	if err != nil {
		return err
	}

	if !caller.IsAdminOfTeam(seat.TeamID) {
		return domain.ErrForbidden
	}

	return nil
}

func (s *CrossTeamRequestService) ApproveRequest(ctx context.Context, id uuid.UUID, callerID uuid.UUID) (*domain.CrossTeamRequest, error) {
	req, err := s.requests.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.canModifyRequest(ctx, callerID, req); err != nil {
		return nil, err
	}

	req.Approve()

	// Auto-reserve: create the reservation when approving
	err = s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.requests.UpdateStatus(txCtx, id, req.Status); err != nil {
			return err
		}

		// Check if already reserved
		existing, _ := s.reservations.GetBySeatAndDate(txCtx, req.TargetSeatID, req.Date)
		if existing == nil {
			reservation := domain.NewReservation(req.RequestingUserID, req.TargetSeatID, req.Date)
			if err := s.reservations.Create(txCtx, reservation); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *CrossTeamRequestService) RejectRequest(ctx context.Context, id uuid.UUID, callerID uuid.UUID) (*domain.CrossTeamRequest, error) {
	req, err := s.requests.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.canModifyRequest(ctx, callerID, req); err != nil {
		return nil, err
	}

	req.Reject()
	if err := s.requests.UpdateStatus(ctx, id, req.Status); err != nil {
		return nil, err
	}
	return req, nil
}
