package ratelimiter

import "errors"

// errors related to rate limiting
var (
	ErrInvalidInterval         = errors.New("Interval must be greater than zero")
	ErrTokenFactoryNotDefined  = errors.New("Token factory must be defined")
	ErrInvalidActiveTokenLimit = errors.New("Active token limit must be greater than zero")
)