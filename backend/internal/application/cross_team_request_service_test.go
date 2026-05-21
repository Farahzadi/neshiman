package application

import (
	"context"
	"errors"
	"neshiman/backend/internal/domain"
	"testing"

	"github.com/google/uuid"
)

func TestCrossTeamRequestService_CreateRequest(t *testing.T) {
	userID := uuid.New()
	seatID := uuid.New()
	date := domain.Date{Year: 2026, Month: 6, Day: 15}

	t.Run("success", func(t *testing.T) {
		var saved *domain.CrossTeamRequest
		svc := NewCrossTeamRequestService(&mockCrossTeamRequestRepo{
			createFn: func(_ context.Context, r *domain.CrossTeamRequest) error {
				saved = r
				return nil
			},
		})

		req, err := svc.CreateRequest(context.Background(), userID, seatID, date)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if req.RequestingUserID != userID || req.TargetSeatID != seatID {
			t.Error("request fields mismatch")
		}
		if req.Status != domain.RequestPending {
			t.Errorf("got %v, want pending", req.Status)
		}
		if saved != req {
			t.Error("saved should be same pointer")
		}
	})

	t.Run("repo error", func(t *testing.T) {
		svc := NewCrossTeamRequestService(&mockCrossTeamRequestRepo{
			createFn: func(_ context.Context, _ *domain.CrossTeamRequest) error {
				return errors.New("db error")
			},
		})
		_, err := svc.CreateRequest(context.Background(), userID, seatID, date)
		if err == nil || err.Error() != "db error" {
			t.Errorf("got %v, want db error", err)
		}
	})
}

func TestCrossTeamRequestService_GetRequest(t *testing.T) {
	id := uuid.New()
	svc := NewCrossTeamRequestService(&mockCrossTeamRequestRepo{
		getByIDFn: func(_ context.Context, got uuid.UUID) (*domain.CrossTeamRequest, error) {
			if got != id {
				t.Errorf("got %v, want %v", got, id)
			}
			return &domain.CrossTeamRequest{ID: id}, nil
		},
	})
	req, err := svc.GetRequest(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.ID != id {
		t.Errorf("got %v, want %v", req.ID, id)
	}
}

func TestCrossTeamRequestService_ListByStatus(t *testing.T) {
	svc := NewCrossTeamRequestService(&mockCrossTeamRequestRepo{
		listByStatusFn: func(_ context.Context, status domain.RequestStatus) ([]domain.CrossTeamRequest, error) {
			if status != domain.RequestPending {
				t.Errorf("got %v, want pending", status)
			}
			return []domain.CrossTeamRequest{{Status: domain.RequestPending}}, nil
		},
	})
	reqs, err := svc.ListByStatus(context.Background(), domain.RequestPending)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reqs) != 1 {
		t.Errorf("got %d, want 1", len(reqs))
	}
}

func TestCrossTeamRequestService_ListPendingByTeam(t *testing.T) {
	teamID := uuid.New()
	svc := NewCrossTeamRequestService(&mockCrossTeamRequestRepo{
		listPendingByTeamFn: func(_ context.Context, got uuid.UUID) ([]domain.CrossTeamRequest, error) {
			if got != teamID {
				t.Errorf("got %v, want %v", got, teamID)
			}
			return []domain.CrossTeamRequest{{Status: domain.RequestPending}}, nil
		},
	})
	reqs, err := svc.ListPendingByTeam(context.Background(), teamID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reqs) != 1 {
		t.Errorf("got %d, want 1", len(reqs))
	}
}

func TestCrossTeamRequestService_ApproveRequest(t *testing.T) {
	id := uuid.New()
	svc := NewCrossTeamRequestService(&mockCrossTeamRequestRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.CrossTeamRequest, error) {
			return &domain.CrossTeamRequest{ID: id, Status: domain.RequestPending}, nil
		},
		updateStatusFn: func(_ context.Context, got uuid.UUID, status domain.RequestStatus) error {
			if got != id {
				t.Errorf("got id %v, want %v", got, id)
			}
			if status != domain.RequestApproved {
				t.Errorf("got %v, want approved", status)
			}
			return nil
		},
	})
	req, err := svc.ApproveRequest(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Status != domain.RequestApproved {
		t.Errorf("got %v, want %v", req.Status, domain.RequestApproved)
	}
}

func TestCrossTeamRequestService_RejectRequest(t *testing.T) {
	id := uuid.New()
	svc := NewCrossTeamRequestService(&mockCrossTeamRequestRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.CrossTeamRequest, error) {
			return &domain.CrossTeamRequest{ID: id, Status: domain.RequestPending}, nil
		},
		updateStatusFn: func(_ context.Context, _ uuid.UUID, status domain.RequestStatus) error {
			if status != domain.RequestRejected {
				t.Errorf("got %v, want rejected", status)
			}
			return nil
		},
	})
	req, err := svc.RejectRequest(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Status != domain.RequestRejected {
		t.Errorf("got %v, want %v", req.Status, domain.RequestRejected)
	}
}
