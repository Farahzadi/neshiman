package dto

type CreateRoomRequest struct {
	Name       string `json:"name"`
	GridWidth  int    `json:"grid_width"`
	GridHeight int    `json:"grid_height"`
}

type UpdateRoomRequest struct {
	Name       string `json:"name"`
	GridWidth  int    `json:"grid_width"`
	GridHeight int    `json:"grid_height"`
}

type CreateReservationRequest struct {
	SeatID string `json:"seat_id"`
	Date   string `json:"date"` // YYYY-MM-DD
}

type CreateTeamRequest struct {
	Name string `json:"name"`
}

type CreateSeatRequest struct {
	RoomID   string `json:"room_id"`
	TeamID   string `json:"team_id"`
	Label    string `json:"label"`
	PosX     int    `json:"pos_x"`
	PosY     int    `json:"pos_y"`
	Rotation int    `json:"rotation"`
}

type MoveSeatRequest struct {
	PosX int `json:"pos_x"`
	PosY int `json:"pos_y"`
}

type CreateUserRequest struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	TeamID      *string `json:"team_id"`
	Role        string  `json:"role"`
	WeeklyLimit int     `json:"weekly_limit"`
}

type UpdateWeeklyLimitRequest struct {
	WeeklyLimit int `json:"weekly_limit"`
}

type CreateCrossTeamRequestRequest struct {
	TargetSeatID string `json:"target_seat_id"`
	Date         string `json:"date"` // YYYY-MM-DD
}

type RotateSeatRequest struct {
	Rotation int `json:"rotation"`
}

type BulkSeatItem struct {
	ID       *string `json:"id"`
	TeamID   string  `json:"team_id"`
	Label    string  `json:"label"`
	PosX     int     `json:"pos_x"`
	PosY     int     `json:"pos_y"`
	Rotation int     `json:"rotation"`
}

type BulkSyncSeatsRequest struct {
	Seats []BulkSeatItem `json:"seats"`
}
