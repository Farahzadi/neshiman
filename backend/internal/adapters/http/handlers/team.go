package handlers

import (
	"encoding/json"
	"net/http"

	"neshiman/backend/internal/adapters/http/dto"
	"neshiman/backend/internal/application"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TeamHandler struct {
	teamSvc *application.TeamService
}

func NewTeamHandler(teamSvc *application.TeamService) *TeamHandler {
	return &TeamHandler{teamSvc: teamSvc}
}

func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	team, err := h.teamSvc.CreateTeam(r.Context(), req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, dto.TeamToResponse(team))
}

func (h *TeamHandler) List(w http.ResponseWriter, r *http.Request) {
	teams, err := h.teamSvc.ListTeams(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responses := make([]dto.TeamResponse, len(teams))
	for i, team := range teams {
		responses[i] = dto.TeamToResponse(&team)
	}
	writeJSON(w, http.StatusOK, responses)
}

func (h *TeamHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid team id", http.StatusBadRequest)
		return
	}
	team, err := h.teamSvc.GetTeam(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, dto.TeamToResponse(team))
}

func (h *TeamHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid team id", http.StatusBadRequest)
		return
	}
	if err := h.teamSvc.DeleteTeam(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
