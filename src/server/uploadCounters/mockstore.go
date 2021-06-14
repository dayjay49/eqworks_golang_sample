package counteruploader

import (
	"fmt"
	"sync"
	"time"
)

type MockStore struct {
	sync.Mutex
	// eventHistory map[string]map[string]int
	EventHistory map[string]map[string]int
}

// NewMockStore creates a new mock store
func NewMockStore() *MockStore {
	return &MockStore{
		EventHistory: make(map[string]map[string]int),
	}
}

// 
func (m *MockStore) AddCounter(c *Counter) {
	dt := time.Now()
	currentDateTime := dt.Format("2006-01-02 15:04:05")

	c.Lock()
	eventKey := c.content + ":" + currentDateTime
	// c.eventKey = dt
	views := c.views
	clicks := c.clicks
	c.Unlock()

	m.Lock()
	m.EventHistory[eventKey] = map[string]int{"views": views, "clicks": clicks}
	// fmt.Println(m.eventHistory)
	m.Unlock()
	fmt.Println("Finished uploading counter with values:")
	fmt.Println(eventKey, views, clicks)
	fmt.Println("-----------------------------------------------------")
}