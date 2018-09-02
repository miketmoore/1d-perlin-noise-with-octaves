package noise

import (
	"math"
	"math/rand"
)

const (
	M = 4294967296
	A = 1664525
	C = 1
)

type PRNG struct {
	Z float64
}

func (p *PRNG) Next() float64 {
	p.Z = math.Mod(A*p.Z+C, M)
	return p.Z/M - 0.5
}

func NewPRNG() PRNG {
	return PRNG{
		Z: math.Floor(rand.Float64() * M),
	}
}

func Interpolate(pa, pb, px float64) float64 {
	ft := px * math.Pi
	f := (1 - math.Cos(ft)) * 0.5
	return pa*(1-f) + pb*f
}

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

func GenerateNoise(amp, wl, octaves, divisor, width float64) []Perlin {
	result := []Perlin{}
	for i := 0.0; i < octaves; i++ {
		result = append(result, NewPerlin(amp, wl, width))
		amp = amp / divisor
		wl = wl / divisor
	}
	return result
}

func CombineNoise(pl []Perlin) []float64 {
	result := []float64{}
	for i, total, j := 0, 0.0, 0; i < len(pl[0].Pos); i++ {
		total = 0
		for j = 0; j < len(pl); j++ {
			total += pl[j].Pos[i]
		}
		result = append(result, total)
	}
	return result
}
