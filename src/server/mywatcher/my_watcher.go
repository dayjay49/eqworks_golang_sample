// package counteruploader
package mywatcher

import (
	"time"

	"github.com/dayjay49/ws-product-golang-master/src/server/mycounters"
	"github.com/dayjay49/ws-product-golang-master/src/server/ratelimiter"
)

func NewMyWatcher(counterConf *mycounters.CounterConfig, rateLimiterConf *ratelimiter.RateLimitConfig) (mycounters.CounterInterface, ratelimiter.RateLimitInterface, error) {
	if counterConf.CycleDuration == 0 {
		return nil, nil, ErrInvalidCycleDuration
	}

	if rateLimiterConf.FixedInterval == 0 {
		return nil, nil, ErrInvalidInterval
	}

	if rateLimiterConf.ActiveTokenLimit == 0 {
		return nil, nil, ErrInvalidActiveTokenLimit
	}

	c := mycounters.NewCounterManager(counterConf)
	r := ratelimiter.NewRateLimitManager(rateLimiterConf)

	w := &ratelimiter.FixedWindowInterval{Interval: rateLimiterConf.FixedInterval}

	// override the manager makeToken function
	r.MakeToken = func() *ratelimiter.Token {
		t := ratelimiter.NewToken()
		t.ExpiresAt = w.EndTime
		return t
	}

	await := func(cycleDuration time.Duration) {
		ticker := time.NewTicker(cycleDuration)
		go func() {
			for {
				select {
				case <-ticker.C:
					// Upload m.counter to mockstore every `cycleDuration` seconds
					c.UploadToMockStore()
				case content := <-c.InFromViewChan:
					// Update counter and then send it to viewHandler
					c.UpdateCounter(content)
				case <-c.MsInChan:
					// Send mock store to statsHandler
					c.SendMockStore()
				case <-r.InTokenReqChan:
					r.TryGenerateToken()
				case token := <-r.ReleaseChan:
					r.ReleaseTokenFromActiveTokenMap(token)
				}
			}
		}()
	}

	w.Run(r.ReleaseExpiredTokens)
	await(counterConf.CycleDuration)
	return c, r, nil
}