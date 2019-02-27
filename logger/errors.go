package logger

import "errors"

var (
	ErrNilLogger         = errors.New("logger is not installed yet")
	ErrBlankOperationKey = errors.New("operation key is blank")
	ErrNilTransfer       = errors.New("transfer should not be nil")
)
