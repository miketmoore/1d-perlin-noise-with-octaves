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

// PRNG is a pseudo-random number generator
type PRNG struct {
	Z float64
}

func (p *PRNG) Next() float64 {
	p.Z = math.Mod(A*p.Z+C, M)
	return p.Z/M - 0.5
}

// NewPRNG returns a new PRNG
func NewPRNG() PRNG {
	return PRNG{
		Z: math.Floor(rand.Float64() * M),
	}
}
