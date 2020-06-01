package input

import "github.com/bsladewski/lapis/stream"

// splitter is the concrete implementation of a stream that splits an input by
// allowing it to be read a specified number of times before advancing to the
// next item in the stream.
type splitter struct {
	n     int
	m     int
	value float64
	in    stream.Stream
	err   error
}

// NewSplitterStream returns a stream that can be used to split an input stream
// amonst multiple output streams by repeating each item in the stream n times.
func NewSplitterStream(in stream.Stream, n int) stream.Stream {

	return &splitter{
		in: in,
		n:  n,
	}

}

func (s *splitter) Next() (float64, error) {

	// if we have consumed the current input n times, read the next value
	if s.m <= 0 {
		s.value, s.err = s.in.Next()
		s.m = s.n
	}

	// decrement the counter that tracks how many times to return the current
	// input
	s.m--

	// if the current error is not nil, return it
	if s.err != nil {
		return 0.0, s.err
	}

	// return the current value
	return s.value, nil

}

func (s *splitter) Close() {
	s.in.Close()
}
