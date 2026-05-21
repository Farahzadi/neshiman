package dto

import (
	"neshiman/backend/internal/domain"
	"time"

	"github.com/google/uuid"
)

type RoomResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	GridWidth  int       `json:"grid_width"`
	GridHeight int       `json:"grid_height"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ReservationResponse struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	SeatID uuid.UUID `json:"seat_id"`
	Date   string    `json:"date"`
}

func RoomToResponse(room *domain.Room) RoomResponse {
	return RoomResponse{
		ID:         room.ID,
		Name:       room.Name,
		GridWidth:  room.GridWidth,
		GridHeight: room.GridHeight,
	}
}

type TeamResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func TeamToResponse(team *domain.Team) TeamResponse {
	return TeamResponse{
		ID:   team.ID,
		Name: team.Name,
	}
}

func ReservationToResponse(r *domain.Reservation) ReservationResponse {
	return ReservationResponse{
		ID:     r.ID,
		UserID: r.UserID,
		SeatID: r.SeatID,
		Date:   formatDate(r.Date),
	}
}

func formatDate(d domain.Date) string {
	return time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
}
