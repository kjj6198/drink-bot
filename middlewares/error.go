package middlewares

import "errors"

var (
	// ErrNoPermission returns when user has no permission
	ErrNoPermission = errors.New("you don't have permission.")
)
