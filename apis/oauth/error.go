package oauth

import (
	"fmt"
)

var (
	// ErrNot17User when user is not from 17.media
	ErrNot17User = fmt.Errorf("you're not 17 user")
	// ErrPermissionDenied when permission is not enough
	ErrPermissionDenied = fmt.Errorf("you don't have permission")
	ErrNoIdTokenField   = fmt.Errorf("request body must have `id_token` field")
)

func makeError(err error) map[string]string {
	return map[string]string{
		"message": err.Error(),
	}
}
