package noise

import "math"

func Interpolate(pa, pb, px float64) float64 {
	ft := px * math.Pi
	f := (1 - math.Cos(ft)) * 0.5
	return pa*(1-f) + pb*f
}
