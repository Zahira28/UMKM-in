package apperror

import (
	"fmt"
	"net/http"
)

type AppError struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	RawErr     error       `json:"-"`
}

func (e *AppError) Error() string {
	if e.RawErr != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.RawErr)
	}
	return e.Message
}

func BadRequest(message string, details ...interface{}) *AppError {
	var det interface{}
	if len(details) > 0 {
		det = details[0]
	}
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		Details:    det,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusForbidden,
		Message:    message,
	}
}

func NotFound(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusNotFound,
		Message:    message,
	}
}

func Conflict(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusConflict,
		Message:    message,
	}
}

func Internal(message string, err ...error) *AppError {
	var raw error
	if len(err) > 0 {
		raw = err[0]
	}
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
		RawErr:     raw,
	}
}
