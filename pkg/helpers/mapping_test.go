package helpers

import (
	"testing"
	"math"
)

const float64EqualityThreshold = 1e-9

func TestNumMap(t *testing.T) {
	nm := NumMap(50, 0, 100, 0, 10)
	comparator := math.Abs(5-nm) <= float64EqualityThreshold
	if !comparator {
		t.Fatalf("Expected 5 but got %f (%v)", nm, comparator)
	}
}