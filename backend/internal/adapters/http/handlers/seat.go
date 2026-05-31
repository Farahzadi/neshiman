package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"neshiman/backend/internal/adapters/http/dto"
	"neshiman/backend/internal/adapters/http/middleware"
	"neshiman/backend/internal/application"
	"neshiman/backend/internal/domain"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SeatHandler struct {
	seatSvc *application.SeatService
	userSvc *application.UserService
}

func NewSeatHandler(seatSvc *application.SeatService, userSvc *application.UserService) *SeatHandler {
	return &SeatHandler{seatSvc: seatSvc, userSvc: userSvc}
}

func (h *SeatHandler) enrichSeatResponse(ctx context.Context, seat *domain.Seat) dto.SeatResponse {
	resp := dto.SeatToResponse(seat)
	if seat.AssignedUserID != nil && h.userSvc != nil {
		if user, err := h.userSvc.GetUser(ctx, *seat.AssignedUserID); err == nil {
			resp.AssignedUserName = user.Name
		}
	}
	return resp
}

func (h *SeatHandler) enrichSeatListResponse(ctx context.Context, seats []domain.Seat) []dto.SeatResponse {
	responses := make([]dto.SeatResponse, len(seats))
	for i, seat := range seats {
		responses[i] = h.enrichSeatResponse(ctx, &seat)
	}
	return responses
}

// CreateSeat creates a new seat
// @Summary      Create a seat
// @Tags         Seats
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateSeatRequest true "Seat details"
// @Success      201  {object}  dto.SeatResponse
// @Failure      400  {string}  string
// @Router       /seats [post]
func (h *SeatHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	roomID, err := uuid.Parse(req.RoomID)
	if err != nil {
		http.Error(w, "invalid room_id", http.StatusBadRequest)
		return
	}
	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		http.Error(w, "invalid team_id", http.StatusBadRequest)
		return
	}
	seat, err := h.seatSvc.CreateSeat(r.Context(), roomID, teamID, req.Label, req.PosX, req.PosY)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, h.enrichSeatResponse(r.Context(), seat))
}

// GetSeatByID returns a seat by ID
// @Summary      Get a seat by ID
// @Tags         Seats
// @Produce      json
// @Param        id   path      string  true  "Seat ID"
// @Success      200  {object}  dto.SeatResponse
// @Failure      404  {string}  string
// @Router       /seats/{id} [get]
func (h *SeatHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid seat id", http.StatusBadRequest)
		return
	}
	seat, err := h.seatSvc.GetSeat(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, h.enrichSeatResponse(r.Context(), seat))
}

// ListSeatsByRoom lists seats for a room
// @Summary      List seats by room
// @Tags         Seats
// @Produce      json
// @Param        room_id  query     string  true  "Room ID"
// @Success      200      {array}   dto.SeatResponse
// @Failure      400      {string}  string
// @Router       /seats [get]
func (h *SeatHandler) ListByRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := uuid.Parse(r.URL.Query().Get("room_id"))
	if err != nil {
		http.Error(w, "invalid or missing room_id query parameter", http.StatusBadRequest)
		return
	}
	seats, err := h.seatSvc.ListSeatsByRoom(r.Context(), roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responses := make([]dto.SeatResponse, len(seats))
	for i, seat := range seats {
		responses[i] = dto.SeatToResponse(&seat)
	}
	writeJSON(w, http.StatusOK, h.enrichSeatListResponse(r.Context(), seats))
}

// MoveSeat moves a seat to a new position
// @Summary      Move a seat
// @Tags         Seats
// @Accept       json
// @Produce      json
// @Param        id       path      string              true  "Seat ID"
// @Param        request  body      dto.MoveSeatRequest  true  "New position"
// @Success      200      {object}  dto.SeatResponse
// @Failure      400      {string}  string
// @Router       /seats/{id}/move [put]
func (h *SeatHandler) Move(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid seat id", http.StatusBadRequest)
		return
	}
	var req dto.MoveSeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	seat, err := h.seatSvc.MoveSeat(r.Context(), id, req.PosX, req.PosY)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, h.enrichSeatResponse(r.Context(), seat))
}

// BulkSyncSeats syncs the full list of seats for a room
// @Summary      Bulk sync seats for a room
// @Tags         Rooms
// @Accept       json
// @Produce      json
// @Param        id       path      string                    true  "Room ID"
// @Param        request  body      dto.BulkSyncSeatsRequest   true  "Full seat list"
// @Success      200      {object}  dto.BulkSyncSeatsResponse
// @Failure      400      {string}  string
// @Router       /rooms/{id}/seats [put]
func (h *SeatHandler) BulkSync(w http.ResponseWriter, r *http.Request) {
	roomID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}
	var req dto.BulkSyncSeatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	seats := make([]domain.Seat, len(req.Seats))
	for i, item := range req.Seats {
		var id uuid.UUID
		if item.ID != nil {
			id, err = uuid.Parse(*item.ID)
			if err != nil {
				http.Error(w, "invalid seat id in list", http.StatusBadRequest)
				return
			}
		}
		teamID, err := uuid.Parse(item.TeamID)
		if err != nil {
			http.Error(w, "invalid team_id in seat list", http.StatusBadRequest)
			return
		}
		seats[i] = domain.Seat{
			ID:     id,
			TeamID: teamID,
			Label:  item.Label,
			Position: domain.Position{
				X: item.PosX,
				Y: item.PosY,
			},
		}
	}
	result, err := h.seatSvc.BulkSyncSeats(r.Context(), roomID, seats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, dto.BulkSyncSeatsResponse{
		Seats: h.enrichSeatListResponse(r.Context(), result),
	})
}

// DeleteSeat deletes a seat
// @Summary      Delete a seat
// @Tags         Seats
// @Param        id   path      string  true  "Seat ID"
// @Success      204  {string}  string
// @Router       /seats/{id} [delete]
func (h *SeatHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid seat id", http.StatusBadRequest)
		return
	}
	if err := h.seatSvc.DeleteSeat(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AssignUser permanently assigns a user to a seat (superadmin only)
// @Summary      Assign a user to a seat permanently
// @Tags         Seats
// @Accept       json
// @Produce      json
// @Param        id       path      string                   true  "Seat ID"
// @Param        request  body      dto.AssignSeatRequest     true  "User ID to assign"
// @Success      200      {object}  dto.SeatResponse
// @Failure      400      {string}  string
// @Failure      403      {string}  string
// @Router       /seats/{id}/assign [post]
func (h *SeatHandler) AssignUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid seat id", http.StatusBadRequest)
		return
	}
	var req dto.AssignSeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}
	callerID := middleware.UserIDFromContext(r.Context())
	seat, err := h.seatSvc.AssignUser(r.Context(), callerID, id, userID)
	if err != nil {
		if err == domain.ErrForbidden {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, h.enrichSeatResponse(r.Context(), seat))
}

// UnassignUser removes permanent user assignment from a seat (superadmin only)
// @Summary      Remove permanent user assignment from a seat
// @Tags         Seats
// @Produce      json
// @Param        id   path      string  true  "Seat ID"
// @Success      200  {object}  dto.SeatResponse
// @Failure      400  {string}  string
// @Failure      403  {string}  string
// @Router       /seats/{id}/assign [delete]
func (h *SeatHandler) UnassignUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid seat id", http.StatusBadRequest)
		return
	}
	callerID := middleware.UserIDFromContext(r.Context())
	seat, err := h.seatSvc.UnassignUser(r.Context(), callerID, id)
	if err != nil {
		if err == domain.ErrForbidden {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, h.enrichSeatResponse(r.Context(), seat))
}
