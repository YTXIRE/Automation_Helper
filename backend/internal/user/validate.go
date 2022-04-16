package user

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func (user CreateUserDTO) Validate() error {
	return validation.ValidateStruct(
		&user,
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Username, validation.Required, validation.Length(4, 32)),
		validation.Field(&user.Password, validation.Required, validation.Length(8, 32)),
	)
}
