package stream_test

import (
	"testing"

	"github.com/bsladewski/lapis/stream"
	"github.com/bsladewski/lapis/util"
)

// TestMAStream tests calculations performed by a Moving Average stream.
func TestMAStream(t *testing.T) {

	// define input data
	inputData := []float64{3.0, 4.0, 2.0, 6.0, 4.0, 5.0, 0.0, 1.0}

	// define expected output averages
	expectedOuput := []float64{3.0, 3.5, 3.0, 4.0, 4.0, 5.0, 3.0, 2.0}

	// create a list stream to provide input to the MA stream
	ls := stream.NewListStream(inputData)

	// create the MA stream
	mas := stream.NewMAStream(ls, 3)
	defer mas.Close()

	// assert that expected output matches data from MA stream
	for i, value := range expectedOuput {

		streamValue, err := mas.Next()
		if err != nil {
			t.Fatalf("index %d; err: %v", i, err)
		}

		t.Logf("expected: %v, got: %v", value, streamValue)

		if util.CompareFloat(streamValue, value) != 0 {
			t.Fatalf("index %d; expected %.2f, got %.2f", i, value, streamValue)
		}

	}

	// assert that next results in end of stream error
	if _, err := mas.Next(); err != stream.ErrEndOfStream {
		t.Fatalf("expected end of stream error, got %v", err)
	}

}
