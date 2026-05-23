package domain

import "errors"

var (
	ErrSeatNotFound          = errors.New("seat not found")
	ErrRoomNotFound          = errors.New("room not found")
	ErrTeamNotFound          = errors.New("team not found")
	ErrUserNotFound          = errors.New("user not found")
	ErrReservationNotFound   = errors.New("reservation not found")
	ErrSeatAlreadyReserved   = errors.New("seat already reserved on this date")
	ErrWeeklyLimitExceeded   = errors.New("weekly reservation limit exceeded")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrForbidden             = errors.New("forbidden")
	ErrCrossTeamRequestNotFound = errors.New("cross-team request not found")
	ErrCrossTeamRequestPending  = errors.New("cross-team request already pending")
	ErrInvalidGridSize       = errors.New("grid width and height must be positive")
	ErrInvalidRotation       = errors.New("rotation must be 0, 90, 180, or 270")
	ErrPositionOutOfBounds   = errors.New("seat position is outside the room grid")
	ErrTeamHasReferences     = errors.New("team has seats or users assigned and cannot be deleted")
)
