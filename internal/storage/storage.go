package storage

import "errors"

var (
	ErrNotExists    = errors.New("not found")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
)
