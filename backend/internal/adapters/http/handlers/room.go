package handlers

import (
	"encoding/json"
	"net/http"

	"neshiman/backend/internal/adapters/http/dto"
	"neshiman/backend/internal/application"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type RoomHandler struct {
	roomSvc *application.RoomService
}

func NewRoomHandler(roomSvc *application.RoomService) *RoomHandler {
	return &RoomHandler{roomSvc: roomSvc}
}

func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	room, err := h.roomSvc.CreateRoom(r.Context(), req.Name, req.GridWidth, req.GridHeight)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, dto.RoomToResponse(room))
}

func (h *RoomHandler) List(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.roomSvc.ListRooms(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responses := make([]dto.RoomResponse, len(rooms))
	for i, room := range rooms {
		responses[i] = dto.RoomToResponse(&room)
	}
	writeJSON(w, http.StatusOK, responses)
}

func (h *RoomHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}
	room, err := h.roomSvc.GetRoom(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, dto.RoomToResponse(room))
}

func (h *RoomHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}
	var req dto.UpdateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	room, err := h.roomSvc.UpdateRoom(r.Context(), id, req.Name, req.GridWidth, req.GridHeight)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, dto.RoomToResponse(room))
}

func (h *RoomHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}
	if err := h.roomSvc.DeleteRoom(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
