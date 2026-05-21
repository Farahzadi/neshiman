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

type SeatResponse struct {
	ID       uuid.UUID `json:"id"`
	RoomID   uuid.UUID `json:"room_id"`
	TeamID   uuid.UUID `json:"team_id"`
	Label    string    `json:"label"`
	PosX     int       `json:"pos_x"`
	PosY     int       `json:"pos_y"`
	Rotation int       `json:"rotation"`
}

func SeatToResponse(seat *domain.Seat) SeatResponse {
	return SeatResponse{
		ID:       seat.ID,
		RoomID:   seat.RoomID,
		TeamID:   seat.TeamID,
		Label:    seat.Label,
		PosX:     seat.Position.X,
		PosY:     seat.Position.Y,
		Rotation: int(seat.Rotation),
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

type UserResponse struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	TeamID      *uuid.UUID `json:"team_id"`
	Role        string     `json:"role"`
	WeeklyLimit int        `json:"weekly_limit"`
}

func UserToResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		TeamID:      user.TeamID,
		Role:        string(user.Role),
		WeeklyLimit: int(user.WeeklyLimit),
	}
}

type CrossTeamRequestResponse struct {
	ID               uuid.UUID `json:"id"`
	RequestingUserID uuid.UUID `json:"requesting_user_id"`
	TargetSeatID     uuid.UUID `json:"target_seat_id"`
	Date             string    `json:"date"`
	Status           string    `json:"status"`
}

func CrossTeamRequestToResponse(r *domain.CrossTeamRequest) CrossTeamRequestResponse {
	return CrossTeamRequestResponse{
		ID:               r.ID,
		RequestingUserID: r.RequestingUserID,
		TargetSeatID:     r.TargetSeatID,
		Date:             formatDate(r.Date),
		Status:           string(r.Status),
	}
}

func formatDate(d domain.Date) string {
	return time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
}
