package output

import (
	"fmt"

	"github.com/bsladewski/gollections"
	"github.com/bsladewski/lapis/stream"
)

// array is used to retrieve a frame of data from a stream as an array.
type array struct {
	size int
	data gollections.Queue
	in   stream.Stream
}

// NewArrayOutput constructs an output that compiles an array representing a
// frame of stream data.
func NewArrayOutput(in stream.Stream, size int) Output {

	return &array{
		size: size,
		data: gollections.NewLinkedQueue(),
		in:   in,
	}

}

func (a *array) Next() (float64, error) {

	// get the next value from the input stream
	value, err := a.in.Next()
	if err != nil {
		return 0.0, err
	}

	// add the input stream to the current frame of data
	a.data.Add(value)
	if a.size > 0 {
		for a.data.Size() > a.size {
			a.data.PopFirst()
		}
	}

	// return the input value
	return value, nil

}

func (a *array) GetData() (interface{}, error) {

	data := []float64{}

	// convert data to an array of floats
	for _, value := range a.data.ToArray() {
		if floatValue, ok := value.(float64); ok {
			data = append(data, floatValue)
		} else {
			return nil, fmt.Errorf("invalid data: %v, Type: %T", value, value)
		}
	}

	return data, nil

}

func (a *array) Close() {
	a.in.Close()
}
