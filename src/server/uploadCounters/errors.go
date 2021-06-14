package counteruploader

import "errors"

// Errors used throughout the codebase
var (
	ErrInvalidCycleDuration = errors.New("Counter upload cycle duration must be greater than zero")
	ErrCounterFactoryNotDefined = errors.New("Counter factory must be defined")
)