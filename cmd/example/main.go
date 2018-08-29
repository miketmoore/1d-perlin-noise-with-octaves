package main

import (
	"fmt"
	"math"
)

// 1D Perlin noise
// Converted from javascript example here:
// https://codepen.io/Tobsta/post/procedural-generation-part-1-1d-perlin-noise

var seed float64 = 100

func main() {
	x := prng()
	fmt.Printf("Pseudo-random number generator example (seed:%f): %f\n", seed, x)

	a := 100.0
	b := 200.0
	mu := 0.5
	y := interpolate(a, b, mu)
	fmt.Printf("Interpolate between %.2f and %.2f (mu:%.2f): %.2f\n", a, b, mu, y)
}

var M float64 = 4294967296

// a - 1 should be divisible by m's prime factors
var A float64 = 1664525

// c and m should be co-prime
var C float64 = 1

var Z float64 = math.Floor(seed)

func prng() float64 {
	Z = math.Mod(A*Z+C, M)
	return Z / M
}

func interpolate(pa, pb, px float64) float64 {
	ft := px * math.Pi
	f := (1 - math.Cos(ft)) * 0.5
	return pa*(1-f) + pb*f
}
