package models

import (
	"github.com/google/uuid"
)

type User struct {
	Uuid      uuid.UUID `json:"uuid"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

type Profile struct {
	Uuid      uuid.UUID `json:"uuid"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	DOB       string    `json:"birthday"`
	Phone     string    `json:"phone"`
	About     string    `json:"about"`
}

type ProfileDto struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	DOB       string `json:"birthday"`
	Phone     string `json:"phone"`
	About     string `json:"about"`
}
