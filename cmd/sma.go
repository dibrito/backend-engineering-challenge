package cmd

import (
	"time"
)

// output representes an event in the output file.
type output struct {
	Date            time.Time `json:"date"` //2018-12-26 18:11:00
	AvgDeliveryTime float32   `json:"average_delivery_time"`
}

// simpleMovingAverage calculates the sma for a given slice of events and writes in
// an output file. Incoming events will be ordered by timestamp.
func simpleMovingAverage(data []event, window int32) {
	// we want to calculate sma for the translation delivery time over the last X minutes.
	// window is already defined.
	// incoming events are already sorted.

	// For each minute, we should get all translation events that occurred within the window time
	// (including the current minute).

	// e.g. At 18:11:
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
	windows := findAllWindows(data[:1][0], data[len(data)-1:][0], window)

	// here we iterate over all windows(ordered and asc)
	// and will get avg delivery time for events that fit in the given window.
	for k, v := range windows {
		// getAvgDeliveryTimeForWindow is itarating over ALL data, N times, where N is the lenght of windows!
		// BAD DECISION!! but let's make it work, then we make it beautiful! ;D
		// complexity: O(nm).
		avg := getAvgDeliveryTimeForWindow(data, v)
		result[k] = output{
			Date:            k,
			AvgDeliveryTime: avg,
		}
	}

	writeOutput(result)
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
