package handlers

import (
	"encoding/json"
	"net/http"

	"neshiman/backend/internal/adapters/http/dto"
	"neshiman/backend/internal/adapters/http/middleware"
	"neshiman/backend/internal/application"
	"neshiman/backend/internal/domain"

	"github.com/google/uuid"
)

type AuthHandler struct {
	authSvc *application.AuthService
}

func NewAuthHandler(authSvc *application.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// Login authenticates a user and returns a JWT token
// @Summary      Login
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Credentials"
// @Success      200    {object}  dto.LoginResponse
// @Failure      400    {string}  string
// @Failure      401    {string}  string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.authSvc.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, dto.LoginResponse{
		Token: result.Token,
		User:  dto.UserToResponse(result.User),
	})
}

// SetPassword sets or updates a user's password
// @Summary      Set password
// @Tags         Auth
// @Accept       json
// @Param        id      path  string                  true  "User ID"
// @Param        request body dto.SetPasswordRequest true "New password"
// @Success      204     {string}  string
// @Failure      400     {string}  string
// @Failure      401     {string}  string
// @Router       /users/{id}/password [put]
func (h *AuthHandler) SetPassword(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.SetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.authSvc.SetPassword(r.Context(), userID, req.Password); err != nil {
		if err == domain.ErrPasswordRequired {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
