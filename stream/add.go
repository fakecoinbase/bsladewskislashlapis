package stream

import (
	"github.com/bsladewski/lapis/util"
	"github.com/pkg/errors"
)

// add is the concrete implementation of a stream that adds together the output
// of a number of input streams; if any input stream reaches the end of input,
// an add stream will output an end of stream error.
type add struct {
	inputs []Stream
}

// NewAddStream returns a stream that adds the output of multiple input streams.
func NewAddStream(inputs ...Stream) Stream {

	return &add{
		inputs: inputs,
	}

}

func (a *add) Next() (float64, error) {

	// errors returned by input streams
	var errs []error

	// the sum of the input stream ouputs
	var sum float64

	// retrieve next value from all input streams, this should consume from
	// each input stream on every call to Next regardless of errors returned
	// by any given input stream
	for _, is := range a.inputs {
		value, err := is.Next()
		if err != nil {
			errs = append(errs, err)
		} else {
			sum += value
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
	return sum, nil

}

func (a *add) Close() {

	// close all input streams
	for _, input := range a.inputs {
		input.Close()
	}

}
