package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"neshiman/backend/internal/adapters/http/dto"
	"neshiman/backend/internal/application"
	"neshiman/backend/internal/domain"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	userSvc *application.UserService
}

func NewUserHandler(userSvc *application.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// CreateUser creates a new user
// @Summary      Create a user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateUserRequest true "User details"
// @Success      201  {object}  dto.UserResponse
// @Failure      400  {string}  string
// @Router       /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	var teamID *uuid.UUID
	if req.TeamID != nil {
		id, err := uuid.Parse(*req.TeamID)
		if err != nil {
			http.Error(w, "invalid team_id", http.StatusBadRequest)
			return
		}
		teamID = &id
	}
	email := ""
	if req.Email != nil {
		email = *req.Email
	}
	user, err := h.userSvc.CreateUser(r.Context(), req.Name, email, teamID, domain.Role(req.Role), domain.WeeklyLimit(req.WeeklyLimit))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, dto.UserToResponse(user))
}

// GetUserByID returns a user by ID
// @Summary      Get a user by ID
// @Tags         Users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  dto.UserResponse
// @Failure      404  {string}  string
// @Router       /users/{id} [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	user, err := h.userSvc.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, dto.UserToResponse(user))
}

// GetUserByEmail returns a user by email
// @Summary      Get a user by email
// @Tags         Users
// @Produce      json
// @Param        email  query     string  true  "User email"
// @Success      200    {object}  dto.UserResponse
// @Failure      404    {string}  string
// @Router       /users/by-email [get]
func (h *UserHandler) GetByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email query parameter required", http.StatusBadRequest)
		return
	}
	user, err := h.userSvc.GetUserByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, dto.UserToResponse(user))
}

// ListUsersByTeam returns users in a team
// @Summary      List users by team
// @Tags         Users
// @Produce      json
// @Param        team_id  query     string  false  "Team ID (omit for all users)"
// @Success      200      {array}   dto.UserResponse
// @Failure      400      {string}  string
// @Router       /users [get]
func (h *UserHandler) ListByTeam(w http.ResponseWriter, r *http.Request) {
	teamIDStr := r.URL.Query().Get("team_id")
	if teamIDStr == "" {
		users, err := h.userSvc.ListAllUsers(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responses := make([]dto.UserResponse, len(users))
		for i, u := range users {
			responses[i] = dto.UserToResponse(&u)
		}
		writeJSON(w, http.StatusOK, responses)
		return
	}
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		http.Error(w, "invalid team_id query parameter", http.StatusBadRequest)
		return
	}
	users, err := h.userSvc.ListUsersByTeam(r.Context(), teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responses := make([]dto.UserResponse, len(users))
	for i, u := range users {
		responses[i] = dto.UserToResponse(&u)
	}
	writeJSON(w, http.StatusOK, responses)
}

// UpdateUserWeeklyLimit updates a user's weekly reservation limit
// @Summary      Update weekly limit
// @Tags         Users
// @Accept       json
// @Param        id       path      string                      true  "User ID"
// @Param        request  body      dto.UpdateWeeklyLimitRequest  true  "Weekly limit"
// @Success      204      {string}  string
// @Failure      400      {string}  string
// @Router       /users/{id}/weekly-limit [put]
func (h *UserHandler) UpdateWeeklyLimit(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	var req dto.UpdateWeeklyLimitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.userSvc.UpdateWeeklyLimit(r.Context(), id, domain.WeeklyLimit(req.WeeklyLimit)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteUser deletes a user
// @Summary      Delete a user
// @Tags         Users
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  dto.DeleteResponse
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	if err := h.userSvc.DeleteUser(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrCannotDeleteSuperAdmin) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, dto.DeleteResponse{Deleted: true})
}
