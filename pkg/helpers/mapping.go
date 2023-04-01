package helpers

// NumMap maps val from in_min/in_max range to out_min/out_max range
func NumMap(val float64, in_min float64, in_max float64, out_min float64, out_max float64) float64 {
	return (val-in_min)*(out_max-out_min)/(in_max-in_min) + out_min
}
