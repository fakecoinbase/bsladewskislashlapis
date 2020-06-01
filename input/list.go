package input

import (
	"fmt"

	"github.com/bsladewski/gollections"
	"github.com/bsladewski/lapis/stream"
)

// A list represents a stream open on a pre-defined list of values.
type list struct {
	values gollections.Queue
}

// NewListStream returns a stream that reads values from a pre-defined list.
func NewListStream(values []float64) stream.Stream {

	// build linked list of values
	valueList := gollections.NewLinkedQueue()
	for _, value := range values {
		valueList.Add(value)
	}

	return &list{
		values: valueList,
	}

}

func (l *list) Next() (float64, error) {

	// return end of stream error if the list is empty
	if l.values.IsEmpty() {
		return 0.0, stream.ErrEndOfStream
	}

	// retrieve and return the first element in the list
	valueI, err := l.values.PopFirst()
	if err != nil {
		return 0.0, err
	}

	if value, ok := valueI.(float64); ok {
		return value, nil
	}

	return 0.0, fmt.Errorf("received invalid value: %v type: %T",
		valueI, valueI)
}

func (l *list) Close() {
	l.values.Clear()
}
