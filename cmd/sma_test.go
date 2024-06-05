package cmd

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewFIFO(t *testing.T) {
	fifo := NewFIFO()
	require.NotNil(t, fifo)
	require.Len(t, fifo.queue, 0)
}

var e = event{
	Timestamp: customTime{
		Time: time.Now().UTC(),
	},
	Duration: 20,
}

func TestFIFOEnqueue(t *testing.T) {
	fifo := NewFIFO()
	require.NotNil(t, fifo)
	require.Len(t, fifo.queue, 0)

	fifo.Enqueue(e)
	require.Len(t, fifo.queue, 1)

	fifo.Enqueue(e)
	require.Len(t, fifo.queue, 2)
}

func TestFIFODequeueNotEmpty(t *testing.T) {
	fifo := NewFIFO()
	require.NotNil(t, fifo)
	require.Len(t, fifo.queue, 0)

	fifo.Enqueue(e)

	require.Len(t, fifo.queue, 1)
	fifo.Dequeue()
	require.Len(t, fifo.queue, 0)
}

func TestFIFODequeueEmpty(t *testing.T) {
	fifo := NewFIFO()
	require.NotNil(t, fifo)
	require.Len(t, fifo.queue, 0)

	fifo.Dequeue()
	require.Len(t, fifo.queue, 0)
}

// lets benchmark:
// we will compare benchmarks for SMA(wich run over ALL events for each minute window)
// and FIFO SMA, which uses an auxiliary FIFO to contain only events for the given minute.
// We will use the tool: benchstat
// go install golang.org/x/perf/cmd/benchstat@latest
// in the following way:
// benchstat sma.bench fifo.bench
// to check performance improvement.
// For now lets use 100k entries.

func generateEventsArray(t *testing.B, numEntries int) []event {
	starTime := time.Now().UTC()
	events := make([]event, numEntries)

	for i := 0; i < numEntries; i++ {
		e, err := generateRandomEvent(starTime)
		require.NoError(t, err)
		events[i] = e
		// increase event timestamp over time.
		starTime = starTime.Add(time.Minute)
	}
	return events
}

func generateRandomEvent(baseTime time.Time) (event, error) {
	// Increment the timestamp by a random number of minutes
	timestamp := baseTime.Add(time.Duration(rand.Intn(60)+1) * time.Minute)
	tt, err := parseTime(timestamp.Format("2006-01-02 15:04:05.999999"))
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return event{}, err
	}
	return event{
		Timestamp: customTime{tt},
		Duration:  rand.Intn(120) + 1,
	}, nil
}

// Prevent inlining of 'leaf functions' and avoid compiler optimizations.
var result map[time.Time]output

func BenchmarkSMA(b *testing.B) {
	// local sink.
	var r map[time.Time]output
	window := int32(10)
	events := generateEventsArray(b, _100K)
	// It's important to not record any setup that is required to run your benchmark.
	b.ResetTimer()
	// execute code to benchmark here:
	for i := 0; i < b.N; i++ {
		r = SMA(events, window)
	}
	result = r
}

func BenchmarkFIFOSMA(b *testing.B) {
	// local sink.
	var r map[time.Time]output
	events := generateEventsArray(b, _100K)
	window := int32(10)
	// It's important to not record any setup that is required to run your benchmark.
	b.ResetTimer()
	// execute code to benchmark here:
	for i := 0; i < b.N; i++ {
		FIFOSMAMinified(events, window)
	}
	result = r
}

func BenchmarkBuffFIFOSMA(b *testing.B) {
	// local sink.
	var r map[time.Time]output
	events := generateEventsArray(b, _100K)
	window := int32(10)
	// It's important to not record any setup that is required to run your benchmark.
	b.ResetTimer()
	// execute code to benchmark here:
	for i := 0; i < b.N; i++ {
		BuffFIFOSMA(events, window)
	}
	result = r
}
