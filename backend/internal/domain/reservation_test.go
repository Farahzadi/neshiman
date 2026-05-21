package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewReservation(t *testing.T) {
	userID := uuid.New()
	seatID := uuid.New()
	date := Date{Year: 2026, Month: 5, Day: 15}

	r := NewReservation(userID, seatID, date)
	if r.UserID != userID {
		t.Errorf("got user %v, want %v", r.UserID, userID)
	}
	if r.SeatID != seatID {
		t.Errorf("got seat %v, want %v", r.SeatID, seatID)
	}
	if r.Date != date {
		t.Errorf("got date %+v, want %+v", r.Date, date)
	}
	if r.ID == uuid.Nil {
		t.Error("expected non-nil ID")
	}
}
