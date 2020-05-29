package stream_test

import (
	"testing"

	"github.com/bsladewski/lapis/stream"
)

// TestAddStream tests adding the result of two streams using an add stream.
func TestAddStream(t *testing.T) {

	// define input data for stream A
	inputDataA := []float64{4.5, 3.9, 8.6, 1.2, 4.4}

	// define input data for stream B
	inputDataB := []float64{9.1, 2.4, 5.2, 2.3}

	// define expected output as sum of streams A and B
	expectedOutput := []float64{13.6, 6.3, 13.8, 3.5}

	// create a list stream for input A
	lsA := stream.NewListStream(inputDataA)

	// create a list stream for input B
	lsB := stream.NewListStream(inputDataB)

	// create the add stream
	as := stream.NewAddStream(lsA, lsB)

	// assert that expected ouput matches data from add stream
	for i, value := range expectedOutput {

		streamValue, err := as.Next()
		if err != nil {
			t.Fatalf("index %d; err: %v", i, err)
		}

		if streamValue != value {
			t.Fatalf("index %d; expected %.2f, got %.2f", i, value, streamValue)
		}

	}

	// assert that next results in end of stream error
	if _, err := as.Next(); err != stream.ErrEndOfStream {
		t.Fatalf("expected end of stream error, got %v", err)
	}

}
