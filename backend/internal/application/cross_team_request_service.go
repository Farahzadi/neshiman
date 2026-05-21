package application

import (
	"context"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"

	"github.com/google/uuid"
)

type CrossTeamRequestService struct {
	requests ports.CrossTeamRequestRepository
}

func NewCrossTeamRequestService(requests ports.CrossTeamRequestRepository) *CrossTeamRequestService {
	return &CrossTeamRequestService{requests: requests}
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

func (s *CrossTeamRequestService) ApproveRequest(ctx context.Context, id uuid.UUID) (*domain.CrossTeamRequest, error) {
	req, err := s.requests.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	req.Approve()
	if err := s.requests.UpdateStatus(ctx, id, req.Status); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *CrossTeamRequestService) RejectRequest(ctx context.Context, id uuid.UUID) (*domain.CrossTeamRequest, error) {
	req, err := s.requests.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	req.Reject()
	if err := s.requests.UpdateStatus(ctx, id, req.Status); err != nil {
		return nil, err
	}
	return req, nil
}
