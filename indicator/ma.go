package indicator

import (
	"errors"
	"fmt"

	"github.com/bsladewski/gollections"
	"github.com/bsladewski/lapis/input"
	"github.com/bsladewski/lapis/math"
	"github.com/bsladewski/lapis/stream"
)

// ma is the concrete implementation of a stream that applies a Moving Average
// function to input data.
type ma struct {
	period int
	frame  gollections.Queue
	in     stream.Stream
}

// NewMAStream returns a stream that applies a Moving Average function to input
// data.
func NewMAStream(in stream.Stream, period int) stream.Stream {

	return &ma{
		period: period,
		frame:  gollections.NewLinkedQueue(),
		in:     in,
	}

}

func (m *ma) Next() (float64, error) {

	if m.period <= 0 {
		return 0.0, errors.New("moving average period cannot be negative or zero")
	}

	// retrieve the next piece of input data
	next, err := m.in.Next()
	if err != nil {
		return 0.0, err
	}

	// add input data to the frame that makes up the current average; if the
	// size of the frame exceeds the period, remove all excess elements
	m.frame.Add(next)
	for m.period > 0 && m.frame.Size() > m.period {
		m.frame.PopFirst()
	}

	// calculate the average of current frame
	sum := 0.0
	for _, valueI := range m.frame.ToArray() {

		if value, ok := valueI.(float64); ok {
			sum += value
		} else {
			return 0.0, fmt.Errorf("received invalid value: %v type: %T",
				valueI, valueI)
		}

	}

	return sum / float64(m.frame.Size()), nil
}

func (m *ma) Close() {
	m.in.Close()
}

// NewMAOscillatorStream returns a stream that applies a Moving Average
// oscillator function to input data.
func NewMAOscillatorStream(in stream.Stream, fastPeriod,
	slowPeriod int) stream.Stream {

	splitter := input.NewSplitterStream(in, 2)

	return math.NewSubStream(
		NewMAStream(splitter, fastPeriod),
		NewMAStream(splitter, slowPeriod),
	)

}
