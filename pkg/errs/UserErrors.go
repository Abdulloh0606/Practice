package errs

import (
	"errors"
)

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")
