package stream_test

import (
	"testing"

	"github.com/bsladewski/lapis/stream"
	"github.com/bsladewski/lapis/util"
)

// TestSplitterStream test the result of splitting an input stream multiple
// ways.
func TestSplitterStream(t *testing.T) {

	// define input data
	inputData := []float64{3.0, 4.0, 2.0, 6.0, 5.0}

	// define expected output values
	expectedOutput := []float64{3.0, 3.0, 3.0, 4.0, 4.0, 4.0, 2.0, 2.0, 2.0,
		6.0, 6.0, 6.0, 5.0, 5.0, 5.0}

	// create a list stream to provide input to the splitter stream
	ls := stream.NewListStream(inputData)

	// create the splitter stream
	ss := stream.NewSplitterStream(ls, 3)
	defer ss.Close()

	// assert that expected ouput matches data from splitter stream
	for i, value := range expectedOutput {

		streamValue, err := ss.Next()
		if err != nil {
			t.Fatalf("index %d; err: %v", i, err)
		}

		t.Logf("expected: %v, got: %v", value, streamValue)

		if util.CompareFloat(streamValue, value) != 0 {
			t.Fatalf("index %d; expected %.2f, got %.2f", i, value, streamValue)
		}

	}

	// assert that next results in end of stream error
	if _, err := ss.Next(); err != stream.ErrEndOfStream {
		t.Fatalf("expected end of stream error, got %v", err)
	}

}
