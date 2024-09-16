package commons

import "errors"

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrInternalServerError = errors.New("internal server error")
	ErrNotFound            = errors.New("not found")
	ErrValidationFailed    = errors.New("validation failed")
	ErrBadRequest          = errors.New("invalid request message")
	ErrTimeout             = errors.New("operation timeout")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrPostNotFound        = errors.New("post not found")
)
