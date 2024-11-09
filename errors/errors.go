package errors

import "errors"

var (
	ErrEmpyVal            = errors.New("empty value")
	ErrEmptyKey           = errors.New("empty key")
	ErrConnClosed         = errors.New("connection closed")
	ErrConnFailed         = errors.New("connection failed")
	ErrEmptyCommand       = errors.New("empty command")
	ErrUnknownCommand     = errors.New("unknown command")
	ErrContentNotFound    = errors.New("content not found")
	ErrGetRequiresOneArg  = errors.New("GET command requires 1 argument")
	ErrDelRequiresOneArg  = errors.New("DEL command requires 1 argument")
	ErrSetRequiresTwoArgs = errors.New("SET command requires 2 arguments")
)
