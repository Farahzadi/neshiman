package handlers

import (
	"encoding/json"
	"net/http"

	"neshiman/backend/internal/adapters/http/dto"
	"neshiman/backend/internal/adapters/http/middleware"
	"neshiman/backend/internal/application"
	"neshiman/backend/internal/domain"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CrossTeamRequestHandler struct {
	svc *application.CrossTeamRequestService
}

func NewCrossTeamRequestHandler(svc *application.CrossTeamRequestService) *CrossTeamRequestHandler {
	return &CrossTeamRequestHandler{svc: svc}
}

func (h *CrossTeamRequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCrossTeamRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	seatID, err := uuid.Parse(req.TargetSeatID)
	if err != nil {
		http.Error(w, "invalid target_seat_id", http.StatusBadRequest)
		return
	}
	date, err := parseDate(req.Date)
	if err != nil {
		http.Error(w, "invalid date, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	request, err := h.svc.CreateRequest(r.Context(), userID, seatID, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, dto.CrossTeamRequestToResponse(request))
}

func (h *CrossTeamRequestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid request id", http.StatusBadRequest)
		return
	}
	request, err := h.svc.GetRequest(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, dto.CrossTeamRequestToResponse(request))
}

func (h *CrossTeamRequestHandler) ListByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		http.Error(w, "status query parameter required", http.StatusBadRequest)
		return
	}
	requests, err := h.svc.ListByStatus(r.Context(), domain.RequestStatus(status))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responses := make([]dto.CrossTeamRequestResponse, len(requests))
	for i, req := range requests {
		responses[i] = dto.CrossTeamRequestToResponse(&req)
	}
	writeJSON(w, http.StatusOK, responses)
}

func (h *CrossTeamRequestHandler) ListPendingByTeam(w http.ResponseWriter, r *http.Request) {
	teamID, err := uuid.Parse(r.URL.Query().Get("team_id"))
	if err != nil {
		http.Error(w, "invalid or missing team_id query parameter", http.StatusBadRequest)
		return
	}
	requests, err := h.svc.ListPendingByTeam(r.Context(), teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responses := make([]dto.CrossTeamRequestResponse, len(requests))
	for i, req := range requests {
		responses[i] = dto.CrossTeamRequestToResponse(&req)
	}
	writeJSON(w, http.StatusOK, responses)
}

func (h *CrossTeamRequestHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid request id", http.StatusBadRequest)
		return
	}
	request, err := h.svc.ApproveRequest(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, dto.CrossTeamRequestToResponse(request))
}

func (h *CrossTeamRequestHandler) Reject(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid request id", http.StatusBadRequest)
		return
	}
	request, err := h.svc.RejectRequest(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, dto.CrossTeamRequestToResponse(request))
}
