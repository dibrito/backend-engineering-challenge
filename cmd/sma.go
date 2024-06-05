package cmd

import (
	"time"
)

// output representes an event in the output file.
type output struct {
	Date            time.Time `json:"date"` //2018-12-26 18:11:00
	AvgDeliveryTime float32   `json:"average_delivery_time"`
}

// SMA calculates the SMA for a given slice of events and writes in
// an output file. Incoming events will be ordered by timestamp.
func SMA(events []event, window int32) map[time.Time]output {
	// we want to calculate sma for the translation delivery time over the last X minutes.
	// window is already defined.
	// incoming events are already sorted.

	// For each minute, we should get all translation events that occurred within the window time
	// (including the current minute).

	// First event at 18:11:
	// window: [18:06, 18:11];
	// Sum the delivery times of all events inside range.
	// Divide the sum by the total number of events.

	// the approach taken will be:
	// we need to know all windows for the incoming events.
	// 1st window(W1) will be: first event time(ET1) minues default window(W): W1 = ET1-W
	// 2nd window(W2) will be: W2 = W1+1min
	// 3rd window(W3) will be: W3 = W2+1min
	// Nth window(WN) will be: WN = W(N-1)+1min.

	result := make(map[time.Time]output)
	// lets find all windows.
	// TODO: rethink this logic of finding all windows, we don't need to find them all,
	// we get the first then iterate/increase them by 1 min.
	windows := findAllWindows(events[:1][0], events[len(events)-1:][0], window)

	// here we iterate over all windows(ordered and asc)
	// and will get avg delivery time for events that fit in the given window.
	for k, v := range windows {
		// getAvgDeliveryTimeForWindow is itarating over ALL data, N times, where N is the lenght of windows!
		// BAD DECISION!! but let's make it work, then we make it beautiful! ;D
		// complexity: O(nm).
		avg := getAvgDeliveryTimeForWindow(events, v)
		result[k] = output{
			Date:            k,
			AvgDeliveryTime: avg,
		}
	}

	return result
}

// getAvgDeliveryTimeForWindow will range over events
// and check for a given pair of window time if the event timestamp
// is between the window range.
func getAvgDeliveryTimeForWindow(events []event, window []time.Time) float32 {
	var sum, count float32
	for _, event := range events {
		if event.Timestamp.Time.After(window[0]) && event.Timestamp.Before(window[1]) {
			// event time is between window.
			sum += float32(event.Duration)
			count++
		}
	}
	if count > 0 {
		return sum / count
	}
	return count
}

// find all windows will find all intervals for each given minute.
// check https://github.com/Unbabel/backend-engineering-challenge/issues/30#issuecomment-550997866
// e.g. given the first event at minute: 18:11, for a given window of 10m,
// the interval to be considered will be: [18:01, 18:11].
func findAllWindows(first, last event, window int32) map[time.Time][]time.Time {
	windows := make(map[time.Time][]time.Time)

	// while current min < last min
	// fill the map

	// who is current?
	// is minute from first event.
	current := getMinute(first)
	// while current is before last event.
	for current.Before(last.Timestamp.Time) {

		if _, ok := windows[current]; !ok {
			// we add the pair(boundaries) to the map of windows.
			windows[current] = append(windows[current], getMinuteDiffRange(current, window), current)
			// and increment current by 1 min.
			current = current.Add(time.Minute)
		}
	}
	// add the last event window
	// current = getMinute(last)
	if _, ok := windows[current]; !ok {
		// we add the pair(boundaries) to the map of windows.
		windows[current] = append(windows[current], getMinuteDiffRange(current, window), current)
	}

	return windows
}

func getMinute(v event) time.Time {
	return v.Timestamp.Truncate(time.Minute)
}

func getMinuteDiffRange(v time.Time, window int32) time.Time {
	return v.Add(time.Duration(int64(time.Minute) * int64(-window))).Truncate(time.Minute)
}

// Thoughts:
// The current issue is we running over all WINDOWS, then for each window we run over all events,
// this is the worst case scenario.
// One importan thing to notice based on SMA:https://en.wikipedia.org/wiki/Moving_average is:

// "When calculating the next mean SMA 𝑘,next with the same sampling width 𝑘 the range from 𝑛 − 𝑘 + 2  to 𝑛 + 1 is considered.
// A new value 𝑝 𝑛 + 1 comes into the sum and the oldest value 𝑝 𝑛 − 𝑘 + 1 drops out.
// This simplifies the calculations by reusing the previous mean SMA 𝑘 , prev."

// In resume: given our first event at: 18:11
// the SMA for the last 5min will be:
// SMA 11 = 10m+9m+8m+7m+6m / total in all minutes
// SMA 12 = 11m+10m+9m+8m+7m / total in all minutes
// SMA 13 = 12m+11m+10m+9m+8m / total in all minutes

// when we say e.g. 6m,xm considere events that happened in the 6 min range or actually before!

// we are always dropping the last "event on given minute",
// as the article says:
// "This means that the moving average filter can be computed quite cheaply on real time data with a FIFO".
// We need to fill the FIFO compute its average on total of elements, assign to a result map, drop last, and move on,
// sounds simple, but let see:

// NOTE: code from here would be 'out of order' since we should have FIFO defined on top of the file
// on in other file, just to show the evolution of 'thought process'!

// lets define out FIFO type:
// FIFO in Go can be obtained with a simple slice
// to enqueue we append, to dequeue we slice of the first element.
type FIFO struct {
	queue []event
}

func NewFIFO() *FIFO {
	return &FIFO{make([]event, 0)}
}

func (f *FIFO) Enqueue(item event) {
	f.queue = append(f.queue, item)
}

func (f *FIFO) Dequeue() {
	if len(f.queue) == 0 {
		return
	}
	f.queue = f.queue[1:]
}

// FIFOSMA calculates sma using FIFO to hold events and avoid iterating over all events.
func FIFOSMA(events []event, window int32) map[time.Time]output {
	fifo := NewFIFO()

	result := make(map[time.Time]output)

	// we want SMA for minute
	// identify range of minutes
	// iterate over minutes
	// fill the fifo, calculate and go to the next minute.

	currMinute := getMinute(events[:1][0])
	end := getMinute(events[len(events)-1:][0])

	// SKIP this comment now:
	// this will mark in each event we currently are!
	// The whole idea here is: we can't/never iterate over events otherwise we'll get worst case scenario!
	currEventIndex := 0

	// For each minute
	// until we reach last minute + 1min, because on the last minute: 18:23:11
	// will still need to calculate and will fall into minute 24!
	for currMinute.Before(end.Add(time.Minute)) || currMinute.Equal(end.Add(time.Minute)) {
		// there will be 3 operations:
		// ADD events,
		// Remove events,
		// Calculate sma.

		// ADD EVENTS:
		// we will add events to fifo that fit in the current time window.

		// The problem with "add events to fifo" is we would need to iterate over all events to see if they "fit"
		// we don't want to iterate over all events.
		// so what we do? We iterate over events AS we ADD them to the fifo!

		// When the event time is over the current minute we say "this is the last event added" and we know it's index.
		// That's how we keep the fifo only with events that fit time
		// and we don't iterate over all events! THIS IS THE REAL "CAT JUMP"!!
		// While the event at the current index falls within the current minute, add it to the queue and move to the next event.
		// ====================================================================================================================
		// ====================================================================================================================

		// REMOVE EVENTS:
		// Remove events from the queue that are older than 10 minutes from the current minute.

		// ====================================================================================================================
		// ====================================================================================================================

		// CALCULATE SMA: for all elements in fifo.
		// ====================================================================================================================
		// ====================================================================================================================

		// ADD EVENTS: in FIFO
		// while the event timestamp is before current minute.
		// beawaer we can't truncate event time here to minutes: cause 18:11:06 insn't in the 18:11 range but in 18:12
		// but currMinute is already truncated to minute!
		// also can never be bigger the events lenght.
		for currEventIndex < len(events) && events[currEventIndex].Timestamp.Before(currMinute) {
			// Add it to the fifo and move to the next event.
			fifo.Enqueue(events[currEventIndex])
			currEventIndex++
		}
		// after this, fifo will have only events that fit into the current minute.

		// remove events that are older than the time window:10 min
		// this is the DROP step where the FIFO will only contain events which fits
		// into [-10min :current event min) ~ [a:b) which means, include elements from index a through b, but not including b
		// similar to slice syntax!
		fifo.queue = dequeueByTime(currMinute, fifo, window)

		// calculate sma for current fifo.
		avg := calculateAvg(fifo)

		// add to result map.
		result[currMinute] = output{
			Date:            currMinute,
			AvgDeliveryTime: avg,
		}

		// increase minute.
		currMinute = currMinute.Add(time.Minute)
	}

	return result
}

// FIFOSMA without comments to make profiling visibility better to understand.
func FIFOSMAMinified(events []event, window int32) map[time.Time]output {
	fifo := NewFIFO()
	result := make(map[time.Time]output)
	currMinute := getMinute(events[:1][0])
	end := getMinute(events[len(events)-1:][0])

	currEventIndex := 0

	for currMinute.Before(end.Add(time.Minute)) || currMinute.Equal(end.Add(time.Minute)) {
		for currEventIndex < len(events) && events[currEventIndex].Timestamp.Before(currMinute) {
			fifo.Enqueue(events[currEventIndex])
			currEventIndex++
		}
		fifo.queue = dequeueByTime(currMinute, fifo, window)
		avg := calculateAvg(fifo)
		result[currMinute] = output{
			Date:            currMinute,
			AvgDeliveryTime: avg,
		}
		currMinute = currMinute.Add(time.Minute)
	}
	return result
}

func BuffFIFOSMA(events []event, window int32) map[time.Time]output {
	result := make(map[time.Time]output)
	currMinute := getMinute(events[:1][0])
	end := getMinute(events[len(events)-1:][0])
	// we might need to start with some capacity!
	// but never start with ZERO!
	// maybe with one!(This was my initial tought!! That cause really bad performance!)
	// fifo := NewBufFIFO(100)
	// after the cpu and mem profiling we come up with at least half of events!
	// didn't work! lets try 10% 25% of events!
	fifo := NewBufFIFO((25 * len(events)) / 100)

	currEventIndex := 0

	for currMinute.Before(end.Add(time.Minute)) || currMinute.Equal(end.Add(time.Minute)) {
		for currEventIndex < len(events) && events[currEventIndex].Timestamp.Before(currMinute) {
			fifo.Enqueue(events[currEventIndex])
			currEventIndex++
		}

		fifo.dequeueBuffFIFOByTime(currMinute, window)
		avg := calculateAvgFromBuffFIFO(fifo)
		result[currMinute] = output{
			Date:            currMinute,
			AvgDeliveryTime: avg,
		}
		currMinute = currMinute.Add(time.Minute)
	}
	return result
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

func calculateAvgFromBuffFIFO(fifo *BufFIFO) float32 {
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

func (fifo *BufFIFO) dequeueBuffFIFOByTime(currMinute time.Time, window int32) {
	// we cant iterate over fifo.queue and remove, so we iterate over a copy.
	auxQueue := fifo.queue
	for _, event := range auxQueue {
		// Remove events from the queue that are older than X minutes from the current minute.
		if currMinute.Sub(event.Timestamp.Time) > time.Minute*time.Duration(window) {
			// remove event from queue.
			fifo.Dequeue()
		}
	}
}
