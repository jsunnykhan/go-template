package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string // never expose in API responses
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

func (u *User) isAdmin() bool { return u.Role == RoleAdmin }

func (u *User) HasRole(r Role) bool { return u.Role == r }
