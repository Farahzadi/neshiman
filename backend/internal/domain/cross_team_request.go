package domain

import "github.com/google/uuid"

type RequestStatus string

const (
	RequestPending  RequestStatus = "pending"
	RequestApproved RequestStatus = "approved"
	RequestRejected RequestStatus = "rejected"
)

type CrossTeamRequest struct {
	ID               uuid.UUID
	RequestingUserID uuid.UUID
	TargetSeatID     uuid.UUID
	Date             Date
	Status           RequestStatus
}

func NewCrossTeamRequest(requestingUserID, targetSeatID uuid.UUID, date Date) *CrossTeamRequest {
	return &CrossTeamRequest{
		ID:               uuid.New(),
		RequestingUserID: requestingUserID,
		TargetSeatID:     targetSeatID,
		Date:             date,
		Status:           RequestPending,
	}
}

func (r *CrossTeamRequest) Approve() {
	r.Status = RequestApproved
}

func (r *CrossTeamRequest) Reject() {
	r.Status = RequestRejected
}
