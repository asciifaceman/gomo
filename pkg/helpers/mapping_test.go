package helpers

import (
	"math"
	"testing"
)

const float64EqualityThreshold = 1e-9

func TestNumReMap(t *testing.T) {
	nm := NumReMap(50, 0, 100, 0, 10)
	comparator := math.Abs(5-nm) <= float64EqualityThreshold
	if !comparator {
		t.Fatalf("Expected 5 but got %f (%v)", nm, comparator)
	}
}
