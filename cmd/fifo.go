package cmd

// BufFIFO represents a circular FIFO.
type BufFIFO struct {
	queue []event
	head  int
	tail  int
	size  int
	cap   int
}

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

func (f *BufFIFO) Dequeue() {
	if f.size == 0 {
		return
	}
	f.head = (f.head + 1) % f.cap
	f.size--
	return
}
