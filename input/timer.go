package input

import (
	"time"

	"github.com/bsladewski/lapis/stream"
)

// A timer is used to introduce a delay before reading the next item in stream.
type timer struct {
	interval time.Duration
	in       stream.Stream
}

// NewTimerStream returns a stream that introduces the specified delay before
// returning the next value in the stream.
func NewTimerStream(in stream.Stream, interval time.Duration) stream.Stream {

	return &timer{
		interval: interval,
		in:       in,
	}

}

func (t *timer) Next() (float64, error) {

	time.Sleep(t.interval)

	return t.in.Next()

}

func (t *timer) Close() {
	t.in.Close()
}
