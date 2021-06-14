package counteruploader

// Manager implements a counter uploader interface
type UploadManager struct {
	// channels
	errorChan		chan error
	outToViewChan   chan *Counter
	inFromViewChan	chan string
	msInChan 		chan struct{}
	msOutChan   	chan *MockStore

	// related structs
	mockStore   	*MockStore
	counter  		*Counter 
}

// NewUploadManager creates a new counter uploader manager
func NewUploadManager(conf *Config) *UploadManager {
	m := &UploadManager {
		errorChan: make(chan error),
		outToViewChan: make(chan *Counter),
		inFromViewChan: make(chan string),
		msInChan: make(chan struct{}),
		msOutChan: make(chan *MockStore),
		mockStore: NewMockStore(),
		counter: NewCounter(conf.InitialContent),
	}
	return m
}

// GetUpdatedCounter is called from viewHandler to update randomized content value
func (u *UploadManager) GetUpdatedCounter(content string) (*Counter, error) {
	go func() {
		u.inFromViewChan <- content
	}()

	// Await for the updated counter
	select {
	case c := <- u.outToViewChan:
		return c, nil
	case err := <- u.errorChan:
		return nil, err
	}
}

// UpdateCounter is called when we need to modify the content of the counter then send it through channel
func (u *UploadManager) UpdateCounter(content string) {
	// change `counter`'s content value
	u.counter.UpdateContent(content)

	// send counter to outToView channel
	go func() {
		u.outToViewChan <- u.counter
	}()
}

// UploadToMockStore is called when we need to upload the counter every 5 seconds
func (u *UploadManager) UploadToMockStore() {
	u.mockStore.AddCounter(u.counter)
}

// GetMockStore is called from statsHandler to get mockstore of counters
func (u *UploadManager) GetMockStore() (*MockStore, error) {
	go func() {
		u.msInChan <- struct{}{}
	}()

	// Await for the mock store
	select {
	case ms := <- u.msOutChan:
		return ms, nil
	case err := <- u.errorChan:
		return nil, err
	}
}

func (u *UploadManager) SendMockStore() {
	// send mock store to msOutChan
	go func() {
		u.msOutChan <- u.mockStore
	}()
}