package cmd

import (
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
