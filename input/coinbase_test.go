package input_test

import (
	"os"
	"testing"

	"github.com/bsladewski/lapis/input"
	"github.com/bsladewski/lapis/stream"
)

// TestCoinbaseStream tests retrieving prices from the coinbase API.
func TestCoinbaseStream(t *testing.T) {

	// construct the coinbase stream
	cs := input.NewCoinbaseStream()

	// retrieve a price
	rate, err := cs.Next()
	if err != nil {
		t.Fatal(err)
	}

	// assert price is greater than zero (if this check fails because the
	// exchange IS zero we've got bigger things to worry about)
	if rate <= 0.0 {
		t.Fatalf("expected rate greater than zero, got %.2f", rate)
	}

	// log the retrieved rate
	t.Log(rate)

}

// TestCoinbaseMockStream tests loading and retrieving historical prices to mock
// coinbase price data.
func TestGetMockSpotPrice(t *testing.T) {

	mockData, err := os.Open("mock_data.csv")
	if err != nil {
		t.Fatal(err)
	}

	// construct the coinbase mock stream
	ms, err := input.NewCoinbaseMockStream(mockData)
	if err != nil {
		t.Fatal(err)
	}

	for {

		// retrieve the next mock price
		rate, err := ms.Next()
		if err == stream.ErrEndOfStream {
			break
		}

		// assert next item of mock data was read without error
		if err != nil {
			t.Fatal(err)
		}

		// assert price is greater than zero (if this check fails because the
		// exchange IS zero we've got bigger things to worry about)
		if rate <= 0.0 {
			t.Fatalf("expected rate greater than zero, got %.2f", rate)
		}

		// log the retrieved rate
		t.Log(rate)

	}

}

// TestGetHistoricalData tests retrieving historical bitcoin data up to the
// current time.
func TestGetHistoricalData(t *testing.T) {

	historicalData, err := input.GetHistoricalData()
	if err != nil {
		t.Fatal(err)
	}

	if len(historicalData) == 0 {
		t.Fatal("failed to retrieve historical data")
	}

}
