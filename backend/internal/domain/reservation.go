package domain

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	SeatID    uuid.UUID
	Date      Date
	DeletedAt *time.Time
}

func NewReservation(userID, seatID uuid.UUID, date Date) *Reservation {
	return &Reservation{
		ID:     uuid.New(),
		UserID: userID,
		SeatID: seatID,
		Date:   date,
	}
}
