package stream

import (
	"github.com/bsladewski/lapis/util"
	"github.com/pkg/errors"
)

// sub is the concrete implementation of a stream that subtracts the output
// of a number of input streams; input streams are subtracted in the order that
// they are added to the sub stream; if any input stream reaches the end of
// input, an add stream will output an end of stream error.
type sub struct {
	inputs []Stream
}

// NewSubStream returns a stream that subtracts the output of multiple input
// streams.
func NewSubStream(inputs ...Stream) Stream {

	return &sub{
		inputs: inputs,
	}

}

func (s *sub) Next() (float64, error) {

	// errors returned by input streams
	var errs []error

	// the result of subtracting the input stream ouputs
	var result float64

	// firstStream notes whether we are looking at the first input stream; if
	// we are reading from the first stream we should set the result variable
	// rather than subtracting from it
	var firstStream = true

	// retrieve next value from all input streams, this should consume from
	// each input stream on every call to Next regardless of errors returned
	// by any given input stream
	for _, is := range s.inputs {
		value, err := is.Next()
		if err != nil {
			errs = append(errs, err)
		} else if firstStream {
			result = value
			firstStream = false
		} else {
			result -= value
		}
	}

	// if we are at the end of any stream, return an end of stream error
	for _, err := range errs {
		for errors.Cause(err) == ErrEndOfStream {
			return 0.0, ErrEndOfStream
		}
	}

	// handle any other errors returned by input streams
	if err := util.ConcatErrors(errs...); err != nil {
		return 0.0, errors.WithStack(err)
	}

	// return the result of the addition
	return result, nil

}

func (s *sub) Close() {

	// close all input streams
	for _, input := range s.inputs {
		input.Close()
	}

}
