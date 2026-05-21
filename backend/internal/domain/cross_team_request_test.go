package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewCrossTeamRequest(t *testing.T) {
	userID := uuid.New()
	seatID := uuid.New()
	date := Date{Year: 2026, Month: 6, Day: 1}

	r := NewCrossTeamRequest(userID, seatID, date)
	if r.RequestingUserID != userID {
		t.Errorf("got user %v, want %v", r.RequestingUserID, userID)
	}
	if r.TargetSeatID != seatID {
		t.Errorf("got seat %v, want %v", r.TargetSeatID, seatID)
	}
	if r.Status != RequestPending {
		t.Errorf("got status %v, want pending", r.Status)
	}
	if r.ID == uuid.Nil {
		t.Error("expected non-nil ID")
	}
}

func TestCrossTeamRequestApprove(t *testing.T) {
	r := NewCrossTeamRequest(uuid.New(), uuid.New(), Date{2026, 6, 1})
	r.Approve()
	if r.Status != RequestApproved {
		t.Errorf("got %v, want %v", r.Status, RequestApproved)
	}
}

func TestCrossTeamRequestReject(t *testing.T) {
	r := NewCrossTeamRequest(uuid.New(), uuid.New(), Date{2026, 6, 1})
	r.Reject()
	if r.Status != RequestRejected {
		t.Errorf("got %v, want %v", r.Status, RequestRejected)
	}
}
