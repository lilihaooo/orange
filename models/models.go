package models

import (
	"time"
)

type Model struct {
	ID        int64      `gorm:"primary_key"json:"id" form:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
