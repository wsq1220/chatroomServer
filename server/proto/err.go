package proto

import "errors"

var (
	ErrUserNotExist  = errors.New("user not exist!")
	ErrInvalidPasswd = errors.New("incorrect user or password!")
)
