package main

import "errors"

var (
	ErrInvalidArg   = errors.New("invalid args")
	ErrUserNotFound = errors.New("user not found")
)
