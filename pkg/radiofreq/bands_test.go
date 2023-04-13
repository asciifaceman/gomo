package radiofreq

import (
	"math"
	"testing"
)

const float64EqualityThreshold = 1e-9

func f64comparator(expected float64, resulted float64) bool {
	return math.Abs(expected-resulted) <= float64EqualityThreshold
}

var tests map[string]float64 = map[string]float64{
	"n71":  0.6,
	"n262": 47,
	"B12":  0.7,
	"B66":  2.1,
}

func TestFrequencyFromShortname(t *testing.T) {

	for k, v := range tests {
		ret := BandMap.FrequencyFromShortname(k)
		if !f64comparator(v, ret) {
			t.Fatalf("Expected response of [%f] but got [%f] for [%s]", v, ret, k)
		}
	}
}

func TestBandFromShortname(t *testing.T) {

	for k, v := range tests {
		ret := BandMap.BandFromShortname(k)
		if !f64comparator(v, ret.Frequency) {
			t.Fatalf("Expected response of [%f] but got [%f] for [%s]", v, ret.Frequency, k)
		}
	}
}
