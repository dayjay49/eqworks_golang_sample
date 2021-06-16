package ratelimiter

import "sync"

type ActiveTokens struct {
	sync.Mutex
	dict map[string]*Token
}

func NewActiveTokens() *ActiveTokens {
	return &ActiveTokens{
		dict: make(map[string]*Token),
	}
}

func (a *ActiveTokens) GetActiveTokenMap() map[string]*Token {
	return a.dict
}

func (a *ActiveTokens) AddActiveToken(token *Token) {
	a.Lock()
	a.dict[token.ID] = token
	a.Unlock()
}

func (a *ActiveTokens) RemoveActiveToken(token *Token) {
	a.Lock()
	delete(a.dict, token.ID)
	a.Unlock()
}