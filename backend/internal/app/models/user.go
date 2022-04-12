package models

import (
	"backend/pkg/helpers"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/m4rw3r/uuid"
)

type User struct {
	ID          uuid.UUID
	Email       string
	Password    string
	RawPassword string
	Login       string
	CreatedAt   int64
	UpdatedAt   int64
	LastLogin   int64
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.RawPassword, validation.Required, validation.Length(8, 32)),
		validation.Field(&u.Login, validation.Required, validation.Length(3, 30)),
	)
}

func (u *User) BeforeCreate() error {
	if len(u.RawPassword) > 0 {
		password, err := helpers.HashPassword(u.RawPassword)
		if err != nil {
			return err
		}
		u.Password = password
	}
	return nil
}
