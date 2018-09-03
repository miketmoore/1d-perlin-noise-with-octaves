package noise

import "math"

type Perlin struct {
	X, Amp, Wl, Fq, A, B float64
	Prng                 PRNG
	Pos                  []float64
}

func NewPerlin(amp, wl, width float64) Perlin {
	p := Perlin{
		X:    0,
		Amp:  amp,
		Wl:   wl,
		Fq:   1 / wl,
		Prng: NewPRNG(),
		Pos:  []float64{},
	}
	p.A = p.Prng.Next()
	p.B = p.Prng.Next()
	for p.X < width {
		if math.Mod(p.X, p.Wl) == 0 {
			p.A = p.B
			p.B = p.Prng.Next()
			p.Pos = append(p.Pos, p.A*p.Amp)
		} else {
			c := math.Mod(p.X, p.Wl) / p.Wl
			interpolated := Interpolate(p.A, p.B, c)
			foo := interpolated * p.Amp
			p.Pos = append(p.Pos, foo)
		}
		p.X++
	}
	return p
}
