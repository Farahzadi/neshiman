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
