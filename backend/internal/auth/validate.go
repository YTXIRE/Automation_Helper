package auth

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

func (auth DTO) Validate() error {
	return validation.ValidateStruct(
		&auth,
		validation.Field(&auth.Username, validation.Required, validation.Length(4, 32)),
		validation.Field(&auth.Password, validation.Required, validation.Length(8, 32)),
	)
}
