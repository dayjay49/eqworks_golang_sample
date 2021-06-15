package ratelimiter

import (
	"errors"
	"fmt"
	"log"
)

// RateLimitManager implements a rate limit interface
type RateLimitManager struct {
	errorChan chan error
	ReleaseChan chan *Token
	InTokenChan	chan struct{}
	outTokenChan chan *Token
	InReqChan chan bool
	outReqChan chan bool
	neededTokensCount int
	activeTokens map[string]*Token
	activeTokenLimit int
	makeToken	tokenFactory
	isAllowed *MoreReqAllowed
}

// NewUploadManager creates a new counter uploader manager
func NewRateLimitManager(conf *RateLimitConfig) *RateLimitManager {
	r := &RateLimitManager {
		errorChan: make(chan error),
		ReleaseChan: make(chan *Token),
		InTokenChan:	make(chan struct{}),
		outTokenChan: make(chan *Token),
		InReqChan: make(chan bool),
		outReqChan: make(chan bool),
		neededTokensCount: 0,
		activeTokens: make(map[string]*Token),
		activeTokenLimit: conf.ActiveTokenLimit,
		makeToken:	NewToken,
		isAllowed: NewMoreReqAllowed(),
	}
	return r
}

// Acquire is called to acquire a new token
func (r *RateLimitManager) Acquire() (*Token, bool, error) {
	go func() {
		r.InTokenChan <- struct{}{}
	}()

	// Await rate limit token
	select {
	case t := <-r.outTokenChan:
		return t, true, nil
	case b := <-r.outReqChan:
		return nil, b, nil
	case err := <-r.errorChan:
		return nil, false, err
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
}

func (r *RateLimitManager) decNeedToken() {
	r.neededTokensCount--
}

func (r *RateLimitManager) awaitingToken() bool {
	return r.neededTokensCount > 0
}

// Called when a new token is needed.
func (r *RateLimitManager) GenerateToken() {
	// panic if token factory is not defined
	if r.makeToken == nil {
		panic(errors.New("Token factory must be defined"))
	}

	// cannot continue if limit has been reached
	if r.isLimitExceeded() {
		r.incNeedToken()
		r.isAllowed.disallowMoreRequests()
		go func() {
			r.InReqChan <- r.isAllowed.value
		}()
		return
	}

	token := r.makeToken()

	// Add token to active map
	r.activeTokens[token.ID] = token

	// send token to outChan
	go func() {
		r.outTokenChan <- token
	}()
}

func (r *RateLimitManager) isLimitExceeded() bool {
	fmt.Println("The number of active tokens is:", len(r.activeTokens))
	fmt.Println("The limit is:", r.activeTokenLimit)
	return len(r.activeTokens) >= r.activeTokenLimit
}

func (r *RateLimitManager) ReleaseTokenFromActiveTokenMap(token *Token) {
	if token == nil {
		log.Print("unable to release nil token")
		return
	}

	if _, ok := r.activeTokens[token.ID]; !ok {
		log.Printf("unable to release token %s - not in use", token.ID)
		return
	}

	// Delete from map
	delete(r.activeTokens, token.ID)

	// process anything waiting for a rate limit
	if r.awaitingToken() {
		r.decNeedToken()
		go r.GenerateToken()
	}
}

// loops over active tokens and releases any that are expired
// for the FixedWindowRateLimiter Algorithm
func (r *RateLimitManager) ReleaseExpiredTokens() {
	for _, token := range r.activeTokens {
		if token.IsExpired() {
			go func(t *Token) {
				r.ReleaseChan <- t
			}(token)
		}
	}
}

// SendIsAllowedValue sends the `False` value to the handlers 
// whenever we need to stop requests
func (r *RateLimitManager) SendIsAllowedValue(value bool) {
	// send the value to 
	go func() {
		r.outReqChan <- r.isAllowed.value
	}()
}