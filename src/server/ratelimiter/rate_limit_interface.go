package ratelimiter

import "time"

type RateLimitInterface interface {
	// // RateLimitManager methods
	Acquire() (*Token, bool, error)
	Release(*Token)
}

// RateLimitConfig represents a rate limit config 
type RateLimitConfig struct {
	// Limit determines how many rate limit tokens can be active at a time
	ActiveTokenLimit int

	// FixedInterval sets the fixed time window for a Fixed Window Rate Limiter
	FixedInterval time.Duration
}