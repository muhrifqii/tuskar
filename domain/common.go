package domain

import "errors"

type ApiResponse[T any] struct {
	Message    string `json:"message"`
	Data       T      `json:"data,omitempty"`
	StatusCode int    `json:"-"`
}

type ApiErrorResponse struct {
	ApiResponse[map[string]interface{}]
}

func (e ApiErrorResponse) Error() string {
	return e.Message
}

var (
	ErrNotFound           = errors.New("error_not_found")
	ErrInvalidCredentials = errors.New("error_invalid_credentials")
	ErrInvalidToken       = errors.New("error_invalid_token")
)

var Errors = []error{
	ErrNotFound,
	ErrInvalidCredentials,
	ErrInvalidToken,
}

var ErrorMessagesByStatus = map[string]string{
	"error_not_found":           "The requested resource was not found",
	"error_invalid_credentials": "The provided credentials are invalid",
	"error_invalid_token":       "The provided token is invalid",
}
