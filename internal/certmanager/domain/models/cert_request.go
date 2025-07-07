package models

import "time"

type CertRequestState string

const (
	CertStatePending  CertRequestState = "pending"
	CertStateApproved CertRequestState = "approved"
	CertStateFailed   CertRequestState = "failed"
)

type CertRequest struct {
	ID        string
	Domain    string
	Status    CertRequestState
	NotBefore *time.Time
	Attempt   int
	Message   *string
	Requested *time.Time
	Created   time.Time
}
