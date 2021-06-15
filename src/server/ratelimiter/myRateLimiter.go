// package counteruploader
package ratelimiter

// NewMyRateLimiter creates an await function for anything related to rate limiting
func NewMyRateLimiter(conf *RateLimitConfig) (RateLimitInterface, error) {
	if conf.FixedInterval == 0 {
		return nil, ErrInvalidInterval
	}

	if conf.ActiveTokenLimit == 0 {
		return nil, ErrInvalidActiveTokenLimit
	}

	r := NewRateLimitManager(conf)
	w := &FixedWindowInterval{interval: conf.FixedInterval}

	// override the manager makeToken function
	r.makeToken = func() *Token {
		t := NewToken()
		t.expiresAt = w.endTime
		return t
	}

	await := func() {
		go func() {
			for {
				select {
				case <-r.InTokenChan:
					r.GenerateToken()
				case token := <-r.ReleaseChan:
					r.ReleaseTokenFromActiveTokenMap(token)
				case b := <-r.InReqChan:
					r.SendIsAllowedValue(b)
				}
			}
		}()
	}

	w.Run(r.ReleaseExpiredTokens)
	await()
	return r, nil
}