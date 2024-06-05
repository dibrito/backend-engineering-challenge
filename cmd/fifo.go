package cmd

import "time"

// FIFO define our FIFO type.
type FIFO struct {
	queue []event
}

// NewFIFO creates a new FIFO.
func NewFIFO() *FIFO {
	return &FIFO{make([]event, 0)}
}

// Enqueue add an item to FIFO.
func (f *FIFO) Enqueue(item event) {
	f.queue = append(f.queue, item)
}

// Dequeue remove 'head' of FIFO.
func (f *FIFO) Dequeue() {
	if len(f.queue) == 0 {
		return
	}
	f.queue = f.queue[1:]
}

// dequeueByTime is a dequeue process that will happen as long as events inside FIFO
// have timestamp Xmin 'smaller' then the minute that is being considere.
func dequeueByTime(currMinute time.Time, fifo *FIFO, window int32) []event {
	// we cant iterate over fifo.queue and remove, so we iterate over a copy.
	auxQueue := fifo.queue
	for _, event := range auxQueue {
		// Remove events from the queue that are older than X minutes from the current minute.
		if currMinute.Sub(event.Timestamp.Time) > time.Minute*time.Duration(window) {
			// remove event from queue.
			fifo.Dequeue()
		}
	}

	return fifo.queue
}

// calculates avg for all elements in FIFO.
func calculateAvg(fifo *FIFO) float32 {
	var sum float32
	for _, event := range fifo.queue {
		sum += float32(event.Duration)
	}
	// avoid division by zero!
	if len(fifo.queue) > 0 {
		return (sum) / float32(len(fifo.queue))
	}
	return 0
}
