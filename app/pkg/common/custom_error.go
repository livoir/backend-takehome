package common

import (
	"net/http"
)

var (
	ErrInternalServerError = NewCustomError(http.StatusInternalServerError, "An internal server error occurred")
	ErrEmailAlreadyExists  = NewCustomError(http.StatusBadRequest, "Email already exists")
	ErrEmailNotFound       = NewCustomError(http.StatusNotFound, "Email not found")
	ErrInvalidPassword     = NewCustomError(http.StatusBadRequest, "Invalid password")
	ErrUserNotFound        = NewCustomError(http.StatusNotFound, "User not found")
	ErrInvalidParam        = NewCustomError(http.StatusBadRequest, "Invalid parameter")
	ErrPostNotFound        = NewCustomError(http.StatusNotFound, "Post not found")
	ErrUnauthorized        = NewCustomError(http.StatusUnauthorized, "Unauthorized")
	ErrInvalidTokenMethod  = NewCustomError(http.StatusUnauthorized, "Invalid token method")
	ErrInvalidToken        = NewCustomError(http.StatusUnauthorized, "Invalid token")
	ErrPostOwnerMismatch   = NewCustomError(http.StatusForbidden, "Post owner mismatch")
)

type CustomError struct {
	StatusCode int
	Message    string
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewCustomError(statusCode int, message string) *CustomError {
	if message == "" {
		message = http.StatusText(statusCode)
	}
	if statusCode < 100 || statusCode > 599 {
		statusCode = http.StatusInternalServerError
	}
	return &CustomError{StatusCode: statusCode, Message: message}
}
