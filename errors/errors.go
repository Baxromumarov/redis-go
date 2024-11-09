package errors

import "errors"

var (
	ErrEmpyVal           = errors.New("empty value")
	ErrEmptyKey          = errors.New("empty key")
	ErrConnClosed        = errors.New("connection closed")
	ErrConnFailed        = errors.New("connection failed")
	ErrKeyNotFound       = errors.New("key not found")
	ErrEmptyCommand      = errors.New("empty command")
	ErrKeyHasExpired     = errors.New("key has expired")
	ErrUnknownCommand    = errors.New("unknown command")
	ErrNoExpirationSet   = errors.New("no expiration set for this key")
	ErrContentNotFound   = errors.New("content not found")
	ErrGetRequiresOneArg = errors.New("GET command requires 1 argument")
	ErrDelRequiresOneArg = errors.New("DEL command requires 1 argument")
	ErrTTLRequiresOneArg = errors.New("TTL command requires 1 argument")
	ErrInvalidTTLVal     = errors.New("invalid TTL value")
	ErrSETCommand        = errors.New("SET command requires at least two arguments: key and value")
)
