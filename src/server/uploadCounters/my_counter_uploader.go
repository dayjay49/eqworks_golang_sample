package counteruploader

import (
	"time"
)

func NewMyCounterUploader(conf *Config) (CounterUploader, error) {
	if conf.CycleDuration == 0 {
		return nil, ErrInvalidCycleDuration
	}

	m := NewUploadManager(conf)

	await := func(cycleDuration time.Duration) {
		ticker := time.NewTicker(cycleDuration)
		go func() {
			for {
				select {
				case <-ticker.C:
					// Upload m.counter to mockstore every `cycleDuration` seconds
					m.UploadToMockStore()
				case content := <- m.inFromViewChan:
					// Update counter and then send it to viewHandler
					m.UpdateCounter(content)
				case <- m.msInChan:
					// Send mock store to statsHandler
					m.SendMockStore()
				}
			}
		}()
	}

	await(conf.CycleDuration)
	return m, nil
}