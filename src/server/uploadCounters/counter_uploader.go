package counteruploader

import "time"

type CounterUploader interface {
	//Manager methods
	GetUpdatedCounter(content string) (*Counter, error)
	GetMockStore() (*MockStore, error)
}

// Config represents a counter uploader config object
type Config struct {
	// CycleDuration is how often the counters are to be uploaded to the mock store
	CycleDuration time.Duration

	InitialContent string
}