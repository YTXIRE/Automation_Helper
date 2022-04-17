package apperror

import "encoding/json"

var (
	ErrNotFound = NewAppError(nil, "not found", "")
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
}

func NewAppError(err error, message, developerMessage string) *AppError {
	return &AppError{
		Err:              err,
		Message:          message,
		DeveloperMessage: developerMessage,
	}
}

func (e AppError) Error() string {
	return e.Message
}

func (e AppError) Unwrap() error {
	return e.Err
}

func (e AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func systemError(err error) *AppError {
	return NewAppError(err, "internal system error", err.Error())
}
