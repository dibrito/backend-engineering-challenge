package cmd

import "time"

// For the sake demonstration, we won't create bufffifo_test.go file
// since we've chosen the FIFOSMS implementation(not the circular).
// This will decrease test coverage, but we are fine with that!

// Just raised test coverage to 78.7% of statements!

// BufFIFO represents a circular FIFO.
type BufFIFO struct {
	queue []event
	head  int
	tail  int
	size  int
	cap   int
}

// NewBufFIFO creates a new BufFIFO with min capacity of 16.
func NewBufFIFO(capacity int) *BufFIFO {
	if capacity < 16 {
		capacity = 16 //min capacity to avoid frequent resizes.
	}
	return &BufFIFO{
		queue: make([]event, capacity),
		head:  0,
		tail:  0,
		size:  0,
		cap:   capacity,
	}
}

// Enqueue enqueue a new event. Double capacity if 'full'.
func (f *BufFIFO) Enqueue(item event) {
	if f.size == f.cap {
		// Expand the buffer if needed
		newCap := f.cap * 2
		newQueue := make([]event, newCap)
		copy(newQueue, f.queue[f.head:])
		copy(newQueue[f.cap-f.head:], f.queue[:f.tail])
		f.queue = newQueue
		f.head = 0
		f.tail = f.size
		f.cap = newCap
	}

	f.queue[f.tail] = item
	f.tail = (f.tail + 1) % f.cap
	f.size++
}

// Dequeue remove the head element.
func (f *BufFIFO) Dequeue() {
	if f.size == 0 {
		return
	}
	f.head = (f.head + 1) % f.cap
	f.size--
}

// dequeueBuffFIFOByTime dequeue all events that meet the given timewindow.
func (fifo *BufFIFO) dequeueBuffFIFOByTime(currMinute time.Time, window int32) {
	if fifo.size == 0 {
		return
	}

	windowDuration := time.Minute * time.Duration(window)

	for fifo.size > 0 {
		event := fifo.queue[fifo.head]

		// Check if the event is within the window.
		if !event.Timestamp.Time.IsZero() && currMinute.Sub(event.Timestamp.Time) > windowDuration {
			// Move the head forward and decrease the size.
			fifo.head = (fifo.head + 1) % fifo.cap
			fifo.size--
		} else {
			// All remaining events are within the window.
			break
		}
	}
}

// calculateAvgFromBuffFIFO calculate avg from all queued elements.
func calculateAvgFromBuffFIFO(fifo *BufFIFO) float32 {
	var sum float32
	count := 0

	for i := 0; i < fifo.size; i++ {
		index := (fifo.head + i) % fifo.cap
		if !fifo.queue[index].Timestamp.Time.IsZero() {
			sum += float32(fifo.queue[index].Duration)
			count++
		}
	}

	if count > 0 {
		return sum / float32(count)
	}
	return 0
}
