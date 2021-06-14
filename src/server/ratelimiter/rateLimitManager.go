package ratelimiter

import (
	"errors"
	"log"
)

// RateLimitManager implements a rate limit interface
type RateLimitManager struct {
	errorChan chan error
	ReleaseChan chan *Token
	InTokenReqChan	chan struct{}
	outTokenChan chan *Token
	neededTokensCount int
	activeTokens map[string]*Token
	activeTokenLimit int
	MakeToken	tokenFactory
}

// NewUploadManager creates a new counter uploader manager
func NewRateLimitManager(conf *RateLimitConfig) *RateLimitManager {
	r := &RateLimitManager {
		errorChan: make(chan error),
		ReleaseChan: make(chan *Token),
		InTokenReqChan:	make(chan struct{}),
		outTokenChan: make(chan *Token),
		neededTokensCount: 0,
		activeTokens: make(map[string]*Token),
		activeTokenLimit: 0,
		MakeToken:	NewToken,
	}
	return r
}

// Acquire is called to acquire a new token
func (r *RateLimitManager) Acquire() (*Token, error) {
	go func() {
		r.InTokenReqChan <- struct{}{}
	}()

	// Await rate limit token
	select {
	case t := <-r.outTokenChan:
		return t, nil
	case err := <-r.errorChan:
		return nil, err
	}
}

// Release is called to release an active token
func (r *RateLimitManager) Release(t *Token) {
	if t.IsExpired() {
		go func() {
			r.ReleaseChan <- t
		}()
	}
}

func (r *RateLimitManager) incNeedToken() {
	r.neededTokensCount++
	// atomic.AddInt64(&m.neededTokensCount, 1)
}

func (r *RateLimitManager) decNeedToken() {
	r.neededTokensCount--
	// atomic.AddInt64(&m.neededTokensCount, -1)
}

func (r *RateLimitManager) awaitingToken() bool {
	// return atomic.LoadInt64(&m.neededTokensCount) > 0
	return r.neededTokensCount > 0
}

// Called when a new token is needed.
func (r *RateLimitManager) TryGenerateToken() {
	// panic if token factory is not defined
	if r.MakeToken == nil {
		panic(errors.New("Token factory must be defined"))
	}

	// cannot continue if limit has been reached
	if r.isLimitExceeded() {
		r.incNeedToken()
		return
	}

	token := r.MakeToken()

	// Add token to active map
	r.activeTokens[token.ID] = token

	// send token to outChan
	go func() {
		r.outTokenChan <- token
	}()
}

func (r *RateLimitManager) isLimitExceeded() bool {
	return len(r.activeTokens) >= r.activeTokenLimit
}

func (r *RateLimitManager) ReleaseTokenFromActiveTokenMap(token *Token) {
	if token == nil {
		log.Print("unable to release nil token")
		return
	}

	if _, ok := r.activeTokens[token.ID]; !ok {
		log.Printf("unable to release token %s - not in use", token)
		return
	}

	// Delete from map
	delete(r.activeTokens, token.ID)

	// process anything waiting for a rate limit
	if r.awaitingToken() {
		r.decNeedToken()
		go r.TryGenerateToken()
	}
}

// loops over active tokens and releases any that are expired
// for FixedWindowRateLimiter Algo
func (r *RateLimitManager) ReleaseExpiredTokens() {
	for _, token := range r.activeTokens {
		if token.IsExpired() {
			go func(t *Token) {
				r.ReleaseChan <- t
			}(token)
		}
	}
}