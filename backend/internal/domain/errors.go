package domain

import "errors"

var (
	ErrSeatNotFound             = errors.New("seat not found")
	ErrRoomNotFound             = errors.New("room not found")
	ErrTeamNotFound             = errors.New("team not found")
	ErrUserNotFound             = errors.New("user not found")
	ErrReservationNotFound      = errors.New("reservation not found")
	ErrSeatAlreadyReserved      = errors.New("seat already reserved on this date")
	ErrWeeklyLimitExceeded      = errors.New("weekly reservation limit exceeded")
	ErrUnauthorized             = errors.New("unauthorized")
	ErrForbidden                = errors.New("forbidden")
	ErrCrossTeamRequestNotFound = errors.New("cross-team request not found")
	ErrCrossTeamRequestPending  = errors.New("cross-team request already pending")
	ErrInvalidGridSize          = errors.New("grid width and height must be positive")
	ErrPositionOutOfBounds      = errors.New("seat position is outside the room grid")
	ErrInvalidCredentials       = errors.New("invalid username or password")
	ErrPasswordRequired         = errors.New("password is required")
	ErrCannotDeleteSuperAdmin   = errors.New("cannot delete the superadmin user")
	ErrWeeklyLimitTooHigh       = errors.New("weekly reservation limit cannot exceed 5")
	ErrNotTeamAdmin             = errors.New("user is not a team admin")
	ErrWrongTeam                = errors.New("user is not in your team")
	ErrCannotChangeSuperAdmin   = errors.New("cannot change the superadmin user")
	ErrPastDate                 = errors.New("cannot reserve for a past date")
	ErrCannotModify            = errors.New("cannot modify this resource")
)
