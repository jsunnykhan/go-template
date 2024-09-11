package models

type Status string
type Gender string

const (
	PENDING    Status = "PENDING"
	ACCEPTED   Status = "ACCEPTED"
	REJECTED   Status = "REJECTED"
	BLOCKED    Status = "BLOCKED"
	UNFRIENDED Status = "UNFRIENDED"
)

const (
	Male   Gender = "Male"
	Female Gender = "Female"
	Other  Gender = "Other"
)

type FriendDto struct {
	Id string `validate:"required,lte=255" json:"id"`
}
