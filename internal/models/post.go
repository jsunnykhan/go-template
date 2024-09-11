package models

import (
	"github.com/google/uuid"
	"time"
)

type PostDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreatedBy struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Post struct {
	Uuid          uuid.UUID `json:"uuid"`
	Title         string    `json:"title"`
	DOC           string    `json:"description"`
	CreatedAt     time.Time `json:"CreatedAt"`
	UpdatedAt     time.Time `json:"UpdatedAt"`
	CreatedByUser CreatedBy `json:"CreatedBy"`
	Likes         int       `json:"likes"`
}
