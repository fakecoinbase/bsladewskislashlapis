// Package coinbase_test contains unit tests for interacting with the coinbase
// API.
package coinbase_test

import (
	"testing"

	"github.com/bsladewski/lapis/coinbase"
)

// TestGetSpotPrice tests retrieving a spot price from the coinbase API.
func TestGetSpotPrice(t *testing.T) {

	// construct the coinbase client
	client := coinbase.NewClient()

	// retrieve spot price
	rate, err := client.GetSpotPrice()
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

// TestGetMockSpotPrice test loading and retrieving historical spot prices.
func TestGetMockSpotPrice(t *testing.T) {

	// construct the coinbase mock client
	client, err := coinbase.NewMockClient()
	if err != nil {
		t.Fatal(err)
	}

	for {

		rate, err := client.GetSpotPrice()
		if err == coinbase.ErrEndOfMockData {
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
