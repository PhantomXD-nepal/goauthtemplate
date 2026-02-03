package types

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email , password or refresh token")
	ErrInternalServer     = errors.New("internal server error")
)
