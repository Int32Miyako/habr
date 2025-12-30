package customerrors

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInternalServer     = errors.New("internal server error")
)

// Auth GRPC errors

var (
	ErrGenerateAccessToken = errors.New("error generating access token")
)
