package mycounters

import "time"

type CounterInterface interface {
	//Manager methods
	GetUpdatedCounter(content string) (*Counter, error)
	GetMockStore() (*MockStore, error)
}

// Config represents a counter uploader config object
type CounterConfig struct {
	// CycleDuration is how often the counters are to be uploaded to the mock store
	CycleDuration time.Duration

	InitialContent string
}