// Package stream defines the fundamental stream interface used to process data
// within the lapis system.
package stream

import "errors"

var (
	// ErrEndOfStream indicates that the end of the stream has been reached.
	ErrEndOfStream = errors.New("end of stream")
)

// A Stream provides allows
type Stream interface {
	// Next gets the next number in the stream.
	Next() (float64, error)
	// Close closes any resources the stream is currently reading.
	Close()
}
