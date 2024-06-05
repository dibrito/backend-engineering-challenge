package cmd

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

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

func TestAllSMAs(t *testing.T) {
	tcs := []struct {
		name        string
		callSMSFunc func([]event) map[time.Time]output
		checkResult func(got map[time.Time]output)
	}{
		{
			name: "when SMA called should create result output",
			callSMSFunc: func(events []event) map[time.Time]output {
				return SMA(events, 10)
			},
		},
		{
			name: "when FIFOSMA called should create result output",
			callSMSFunc: func(events []event) map[time.Time]output {
				return FIFOSMAMinified(events, 10)
			},
		},
		{
			name: "when BuffFIFOSMA called should create result output",
			callSMSFunc: func(events []event) map[time.Time]output {
				return BuffFIFOSMA(events, 10)
			},
		},
	}

	events, err := parseInputFile("../events.json")
	require.NoError(t, err)
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			got := tc.callSMSFunc(events)
			want := createWantOutput()

			if len(want) != len(got) {
				t.Fatalf("Expected map length %d, got %d", len(want), len(got))
			}
			for k, v := range want {
				if got[k] != v {
					t.Errorf("At %v: expected %+v, got %+v", k, v, got[k])
				}
			}

		})
	}
}

func createWantOutput() map[time.Time]output {
	layout := "2006-01-02 15:04:05"
	want := make(map[time.Time]output)

	dates := []string{
		"2018-12-26 18:11:00",
		"2018-12-26 18:12:00",
		"2018-12-26 18:13:00",
		"2018-12-26 18:14:00",
		"2018-12-26 18:15:00",
		"2018-12-26 18:16:00",
		"2018-12-26 18:17:00",
		"2018-12-26 18:18:00",
		"2018-12-26 18:19:00",
		"2018-12-26 18:20:00",
		"2018-12-26 18:21:00",
		"2018-12-26 18:22:00",
		"2018-12-26 18:23:00",
		"2018-12-26 18:24:00",
	}

	avgTimes := []float32{
		0,
		20,
		20,
		20,
		20,
		25.5,
		25.5,
		25.5,
		25.5,
		25.5,
		25.5,
		31,
		31,
		42.5,
	}

	for i, dateStr := range dates {
		date, err := time.Parse(layout, dateStr)
		if err != nil {
			panic(err)
		}
		want[date] = output{
			Date:            date,
			AvgDeliveryTime: avgTimes[i],
		}
	}

	return want
}
