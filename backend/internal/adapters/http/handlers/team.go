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

// CreateTeam creates a new team
// @Summary      Create a team
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateTeamRequest true "Team name"
// @Success      201  {object}  dto.TeamResponse
// @Failure      400  {string}  string
// @Router       /teams [post]
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

// ListTeams returns all teams
// @Summary      List teams
// @Tags         Teams
// @Produce      json
// @Success      200  {array}   dto.TeamResponse
// @Router       /teams [get]
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

// GetTeamByID returns a team by ID
// @Summary      Get a team by ID
// @Tags         Teams
// @Produce      json
// @Param        id   path      string  true  "Team ID"
// @Success      200  {object}  dto.TeamResponse
// @Failure      404  {string}  string
// @Router       /teams/{id} [get]
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

// DeleteTeam deletes a team
// @Summary      Delete a team
// @Tags         Teams
// @Param        id   path      string  true  "Team ID"
// @Success      200  {object}  dto.DeleteResponse
// @Failure      400  {string}  string
// @Failure      409  {string}  string
// @Router       /teams/{id} [delete]
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
	writeJSON(w, http.StatusOK, dto.DeleteResponse{Deleted: true})
}
