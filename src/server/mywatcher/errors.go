package mywatcher

import "errors"

// Errors used throughout the codebase
var (
	ErrInvalidCycleDuration = errors.New("Counter upload cycle duration must be greater than zero")
	
	ErrInvalidInterval         = errors.New("Interval must be greater than zero")
	ErrTokenFactoryNotDefined  = errors.New("Token factory must be defined")
	ErrInvalidActiveTokenLimit = errors.New("Active token limit must be greater than zero")
)