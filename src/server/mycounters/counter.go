package mycounters

import (
	"fmt"
	"sync"
	"time"
)

type Counter struct {
	sync.Mutex
	views  int
	clicks int
	content string
	// eventKey time.Time
	// eventValue map[string]int
	createdAt time.Time
	updatedAt time.Time
}

// NewCounter creates a new counter
func NewCounter(data string) *Counter {
	return &Counter{
		views: 0,
		clicks: 0,
		content: data,
		createdAt: time.Now().UTC(),
		updatedAt: time.Now().UTC(),
	}
}

// UpdateContent updates the `content` attribute
func (c *Counter) UpdateContent(content string) {
	c.Lock()
	c.content = content
	c.updatedAt = time.Now().UTC()
	fmt.Println("UPDATED COUNTER CONTENT VALUE TO:", c.content)
	c.Unlock()
}

// IncrementView increments the view count of the counter struct
func (c *Counter) IncrementView() {
	c.Lock()
	c.views++
	c.Unlock()
}

// IncrementClick increments the click count of the counter struct
func (c *Counter) IncrementClick() {
	c.Lock()
	c.clicks++
	c.Unlock()
}