package ratelimiter

import (
	"time"

	"github.com/segmentio/ksuid"
)

// token factory function creates a new token
type tokenFactory func() *Token

// Token represents a Rate Limit Token
type Token struct {
	// The unique token ID
	ID string

	// The time at which the token was created
	CreatedAt time.Time

	// Defines the min amount of time the token must live before being released
	ExpiresAt time.Time
}

// NewToken creates a new token
func NewToken() *Token {
	return &Token{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Time{}, // defaults to zero time
	}
}

// IsExpired returns true if current time is greater than expiration time
func (t *Token) IsExpired() bool {
	now := time.Now().UTC()
	return t.ExpiresAt.Before(now)
}