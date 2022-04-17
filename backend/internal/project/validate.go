package project

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func (project Project) Validate() error {
	return validation.ValidateStruct(
		&project,
		validation.Field(&project.Name, validation.Required, validation.Length(4, 32)),
		validation.Field(&project.TestRailProjectID, validation.Required, validation.Length(1, 100), is.Int),
	)
}
