package mycounters

import "errors"

// Errors used throughout the codebase
var (
	ErrInvalidCycleDuration = errors.New("Counter upload cycle duration must be greater than zero")
)