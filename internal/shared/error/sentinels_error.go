package error

import "errors"

var (
	ErrFailedToCreateLogFiles = errors.New("failed to create log file")
)
