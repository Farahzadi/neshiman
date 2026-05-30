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

// CreateCrossTeamRequest creates a cross-team seat request
// @Summary      Create a cross-team request
// @Tags         CrossTeamRequests
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateCrossTeamRequestRequest true "Request details"
// @Success      201  {object}  dto.CrossTeamRequestResponse
// @Failure      400  {string}  string
// @Router       /cross-team-requests [post]
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

// GetCrossTeamRequestByID returns a cross-team request by ID
// @Summary      Get a cross-team request by ID
// @Tags         CrossTeamRequests
// @Produce      json
// @Param        id   path      string  true  "Request ID"
// @Success      200  {object}  dto.CrossTeamRequestResponse
// @Failure      404  {string}  string
// @Router       /cross-team-requests/{id} [get]
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

// ListCrossTeamRequestsByStatus lists requests by status
// @Summary      List requests by status
// @Tags         CrossTeamRequests
// @Produce      json
// @Param        status  query     string  true  "Status (pending/approved/rejected)"
// @Success      200     {array}   dto.CrossTeamRequestResponse
// @Failure      400     {string}  string
// @Router       /cross-team-requests [get]
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

// ListPendingCrossTeamRequestsByTeam lists pending requests for a team
// @Summary      List pending requests by team
// @Tags         CrossTeamRequests
// @Produce      json
// @Param        team_id  query     string  true  "Team ID"
// @Success      200      {array}   dto.CrossTeamRequestResponse
// @Failure      400      {string}  string
// @Router       /cross-team-requests/pending-by-team [get]
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

// ListMyCrossTeamRequests lists the current user's own requests
// @Summary      List my cross-team requests
// @Tags         CrossTeamRequests
// @Produce      json
// @Success      200  {array}   dto.CrossTeamRequestResponse
// @Router       /cross-team-requests/mine [get]
func (h *CrossTeamRequestHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	requests, err := h.svc.ListByUser(r.Context(), userID)
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

// ListAllCrossTeamRequests lists all requests (superadmin only)
// @Summary      List all cross-team requests
// @Tags         CrossTeamRequests
// @Produce      json
// @Success      200  {array}   dto.CrossTeamRequestResponse
// @Router       /cross-team-requests/all [get]
func (h *CrossTeamRequestHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	requests, err := h.svc.ListAll(r.Context())
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

// ApproveCrossTeamRequest approves a pending request
// @Summary      Approve a cross-team request
// @Tags         CrossTeamRequests
// @Produce      json
// @Param        id   path      string  true  "Request ID"
// @Success      200  {object}  dto.CrossTeamRequestResponse
// @Failure      400  {string}  string
// @Failure      403  {string}  string
// @Router       /cross-team-requests/{id}/approve [put]
func (h *CrossTeamRequestHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid request id", http.StatusBadRequest)
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	request, err := h.svc.ApproveRequest(r.Context(), id, userID)
	if err != nil {
		status := http.StatusBadRequest
		if err == domain.ErrForbidden {
			status = http.StatusForbidden
		}
		http.Error(w, err.Error(), status)
		return
	}
	writeJSON(w, http.StatusOK, dto.CrossTeamRequestToResponse(request))
}

// RejectCrossTeamRequest rejects a pending request
// @Summary      Reject a cross-team request
// @Tags         CrossTeamRequests
// @Produce      json
// @Param        id   path      string  true  "Request ID"
// @Success      200  {object}  dto.CrossTeamRequestResponse
// @Failure      400  {string}  string
// @Failure      403  {string}  string
// @Router       /cross-team-requests/{id}/reject [put]
func (h *CrossTeamRequestHandler) Reject(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid request id", http.StatusBadRequest)
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	request, err := h.svc.RejectRequest(r.Context(), id, userID)
	if err != nil {
		status := http.StatusBadRequest
		if err == domain.ErrForbidden {
			status = http.StatusForbidden
		}
		http.Error(w, err.Error(), status)
		return
	}
	writeJSON(w, http.StatusOK, dto.CrossTeamRequestToResponse(request))
}
