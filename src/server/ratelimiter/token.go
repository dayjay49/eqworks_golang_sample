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
	createdAt time.Time

	// Defines the min amount of time the token must live before being released
	expiresAt time.Time
}

// NewToken creates a new token
func NewToken() *Token {
	return &Token{
		ID:        ksuid.New().String(),
		createdAt: time.Now().UTC(),
		expiresAt: time.Time{}, // defaults to zero time
	}
}

// IsExpired returns true if current time is greater than expiration time
func (t *Token) IsExpired() bool {
	now := time.Now().UTC()
	return t.expiresAt.Before(now)
}