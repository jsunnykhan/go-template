package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Name struct {
	gorm.Model
	Uuid         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();" json:"uuid"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
}

