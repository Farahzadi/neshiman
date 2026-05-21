package handlers

import (
	"encoding/json"
	"net/http"

	"neshiman/backend/internal/adapters/http/dto"
	"neshiman/backend/internal/application"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SeatHandler struct {
	seatSvc *application.SeatService
}

func NewSeatHandler(seatSvc *application.SeatService) *SeatHandler {
	return &SeatHandler{seatSvc: seatSvc}
}

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
	seat, err := h.seatSvc.CreateSeat(r.Context(), roomID, teamID, req.Label, req.PosX, req.PosY, req.Rotation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, dto.SeatToResponse(seat))
}

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
	writeJSON(w, http.StatusOK, dto.SeatToResponse(seat))
}

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
	writeJSON(w, http.StatusOK, responses)
}

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
	writeJSON(w, http.StatusOK, dto.SeatToResponse(seat))
}

func (h *SeatHandler) Rotate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid seat id", http.StatusBadRequest)
		return
	}
	var req dto.RotateSeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	seat, err := h.seatSvc.RotateSeat(r.Context(), id, req.Rotation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, dto.SeatToResponse(seat))
}

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
