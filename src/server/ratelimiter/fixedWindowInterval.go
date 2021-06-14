package ratelimiter

import "time"

// FixedWindowInterval represents a fixed window of time with a start / end time
type FixedWindowInterval struct {
	startTime time.Time
	EndTime   time.Time
	Interval  time.Duration
}

func (w *FixedWindowInterval) setWindowTime() {
	w.startTime = time.Now().UTC()
	w.EndTime = time.Now().UTC().Add(w.Interval)
}

func (w *FixedWindowInterval) Run(cb func()) {
	go func() {
		ticker := time.NewTicker(w.Interval)
		w.setWindowTime()
		for range ticker.C {
			cb()
			w.setWindowTime()
		}
	}()
}