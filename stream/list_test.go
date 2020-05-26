package stream_test

import (
	"testing"

	"github.com/bsladewski/lapis/stream"
)

// TestListStream tests reading values from a list stream.
func TestListStream(t *testing.T) {

	// define input data
	inputData := []float64{1.0, 2.1, 3.2, 4.3, 5.4, 6.5}

	// create the stream
	ls := stream.NewListStream(inputData)
	defer ls.Close()

	// assert that input data matches data from stream
	for i, value := range inputData {

		streamValue, err := ls.Next()
		if err != nil {
			t.Fatalf("index %d; err: %v", i, err)
		}

		if streamValue != value {
			t.Fatalf("index %d; expected %.2f, got %.2f", i, value, streamValue)
		}

	}

	// assert that next results in end of stream error
	if _, err := ls.Next(); err != stream.ErrEndOfStream {
		t.Fatalf("expected end of stream error, got %v", err)
	}

}
