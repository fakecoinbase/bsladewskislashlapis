// Package coinbase_test contains unit tests for interacting with the coinbase
// API.
package coinbase_test

import (
	"os"
	"testing"

	"github.com/bsladewski/lapis/coinbase"
	"github.com/bsladewski/lapis/stream"
)

// TestGetSpotPrice tests retrieving a spot price from the coinbase API.
func TestGetSpotPrice(t *testing.T) {

	// construct the coinbase client
	client := coinbase.NewClient()

	// retrieve spot price
	rate, err := client.Next()
	if err != nil {
		t.Fatal(err)
	}

	// assert spot price is greater than zero (if this check fails because the
	// exchange IS zero we've got bigger things to worry about)
	if rate <= 0.0 {
		t.Fatalf("expected rate greater than zero, got %.2f", rate)
	}

	// log the retrieved rate
	t.Log(rate)

}

// TestGetHistoricalData tests retrieving historical bitcoin data up to the
// current time.
func TestGetHistoricalData(t *testing.T) {

	historicalData, err := coinbase.GetHistoricalData()
	if err != nil {
		t.Fatal(err)
	}

	if len(historicalData) == 0 {
		t.Fatal("failed to retrieve historical data")
	}

}

// TestGetMockSpotPrice test loading and retrieving historical spot prices to
// mock coinbase spot price data.
func TestGetMockSpotPrice(t *testing.T) {

	mockData, err := os.Open("mock_data.csv")
	if err != nil {
		t.Fatal(err)
	}

	// construct the coinbase mock client
	client, err := coinbase.NewMockClient(mockData)
	if err != nil {
		t.Fatal(err)
	}

	for {

		rate, err := client.Next()
		if err == stream.ErrEndOfStream {
			break
		}

		// assert next item of mock data was read without error
		if err != nil {
			t.Fatal(err)
		}

		// assert spot price is greater than zero (if this check fails because the
		// exchange IS zero we've got bigger things to worry about)
		if rate <= 0.0 {
			t.Fatalf("expected rate greater than zero, got %.2f", rate)
		}

		// log the retrieved rate
		t.Log(rate)

	}

}
