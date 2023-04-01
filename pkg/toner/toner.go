package toner

import "math/rand"

type Tone struct {
}

func (t *Tone) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		samples[i][0] = rand.Float64()*2 - 1
		samples[i][1] = rand.Float64()*2 - 1
	}
	return len(samples), true
}

func (t *Tone) Err() error {
	return nil
}
