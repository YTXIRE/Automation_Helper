package models

import (
	"github.com/m4rw3r/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Password  string
	Login     string
	CreatedAt int64
	UpdatedAt int64
	LastLogin int64
}
