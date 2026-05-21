package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"neshiman/backend/internal/adapters/http/dto"
	"neshiman/backend/internal/adapters/http/middleware"
	"neshiman/backend/internal/application"
	"neshiman/backend/internal/domain"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ReservationHandler struct {
	reservationSvc *application.ReservationService
}

func NewReservationHandler(reservationSvc *application.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationSvc: reservationSvc}
}

func (h *ReservationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	seatID, err := uuid.Parse(req.SeatID)
	if err != nil {
		http.Error(w, "invalid seat id", http.StatusBadRequest)
		return
	}
	date, err := parseDate(req.Date)
	if err != nil {
		http.Error(w, "invalid date, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	reservation, err := h.reservationSvc.ReserveSeat(r.Context(), userID, seatID, date)
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case domain.ErrSeatAlreadyReserved:
			status = http.StatusConflict
		case domain.ErrWeeklyLimitExceeded:
			status = http.StatusTooManyRequests
		case domain.ErrSeatNotFound, domain.ErrUserNotFound:
			status = http.StatusNotFound
		case domain.ErrForbidden:
			status = http.StatusForbidden
		}
		http.Error(w, err.Error(), status)
		return
	}
	writeJSON(w, http.StatusCreated, dto.ReservationToResponse(reservation))
}

func (h *ReservationHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	reservationID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid reservation id", http.StatusBadRequest)
		return
	}
	if err := h.reservationSvc.CancelReservation(r.Context(), reservationID, userID); err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrForbidden {
			status = http.StatusForbidden
		} else if err == domain.ErrReservationNotFound {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ReservationHandler) GetWeek(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	dateStr := r.URL.Query().Get("date")
	date, err := parseDate(dateStr)
	if err != nil {
		date = domain.Date{Year: time.Now().Year(), Month: int(time.Now().Month()), Day: time.Now().Day()}
	}
	reservations, err := h.reservationSvc.GetUserWeekReservations(r.Context(), userID, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responses := make([]dto.ReservationResponse, len(reservations))
	for i, r := range reservations {
		responses[i] = dto.ReservationToResponse(&r)
	}
	writeJSON(w, http.StatusOK, responses)
}

func parseDate(s string) (domain.Date, error) {
	if s == "" {
		t := time.Now()
		return domain.Date{Year: t.Year(), Month: int(t.Month()), Day: t.Day()}, nil
	}
	t, err := time.Parse("2006-01-02", strings.TrimSpace(s))
	if err != nil {
		return domain.Date{}, err
	}
	return domain.Date{Year: t.Year(), Month: int(t.Month()), Day: t.Day()}, nil
}
