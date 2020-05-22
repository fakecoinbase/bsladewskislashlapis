// Package coinbase_test contains unit tests for interacting with the coinbase
// API.
package coinbase_test

import (
	"testing"

	"github.com/bsladewski/lapis/coinbase"
)

// TestGetSpotPrice tests retrieving a spot price from the coinbase API.
func TestGetSpotPrice(t *testing.T) {

	// retrieve spot price
	rate, err := coinbase.GetSpotPrice()
	if err != nil {
		t.Fatal(err)
	}

	// assert spot price is greater than zero (if this check fails because the
	// exchange IS zero we've got bigger things to worry about)
	if rate <= 0.0 {
		t.Fatalf("expected rate greater than zero, got %.2f", rate)
	}

	// log the retrieve rate
	t.Log(rate)

}
