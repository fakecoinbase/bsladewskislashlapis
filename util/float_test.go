package util_test

import (
	"testing"

	"github.com/bsladewski/lapis/util"
)

// TestCompareFloat tests the results of CompareFloat function under a variety
// of cases.
func TestCompareFloat(t *testing.T) {

	// define test cases
	cases := []struct {
		name     string
		a        float64
		b        float64
		expected int
	}{
		{"TestEqual", -1.1, -1.0999999999, 0},
		{"TestLessThan", -1.1, -1.099999999, -1},
		{"TestGreaterThan", -1.1, -1.1000000001, 0},
	}

	// run each test case
	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {

			// assert that the CompareFloat function returns the expected output
			if got := util.CompareFloat(tc.a, tc.b); got != tc.expected {
				t.Fatalf("CompareFloat(%v, %v), expected %d, got %d", tc.a,
					tc.b, tc.expected, got)
			}

		})

	}

}
