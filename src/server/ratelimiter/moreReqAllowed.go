package ratelimiter

import (
	"fmt"
	"sync"
)

// MoreReqAllowed represents whether more requests should be allowed or not
type MoreReqAllowed struct {
	sync.Mutex
	value bool
}

func NewMoreReqAllowed() *MoreReqAllowed {
	return &MoreReqAllowed{
		value: true,
	}
}

func (a *MoreReqAllowed) disallowMoreRequests() {
	a.Lock()
	a.value = false
	a.Unlock()
	fmt.Println("-----DISALLOWING ANY MORE REQUESTS UNTIL NEXT WINDOW-----")
}