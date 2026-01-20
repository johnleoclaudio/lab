package models

import "errors"

var (
	ErrNotFound           = errors.New("resource not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrValidation         = errors.New("validation failed")
	ErrConflict           = errors.New("resource conflict")
	ErrEmailAlreadyExists = errors.New("email already exists")
)
