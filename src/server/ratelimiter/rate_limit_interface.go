package ratelimiter

import "time"

type RateLimitInterface interface {
	// // RateLimitManager methods
	Acquire() (*Token, error)
	Release(*Token)
}

// RateLimitConfig represents a rate limit config 
type RateLimitConfig struct {
	ActiveTokenLimit int

	FixedInterval time.Duration
}