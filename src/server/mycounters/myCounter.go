package mycounters

import (
	"time"
)

// NewMyCounter creates an await function for anything related to counters and mock store
func NewMyCounter(conf *CounterConfig) (CounterInterface, error) {
	if conf.CycleDuration == 0 {
		return nil, ErrInvalidCycleDuration
	}

	c := NewCounterManager(conf)

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
				}
			}
		}()
	}

	await(conf.CycleDuration)
	return c, nil
}