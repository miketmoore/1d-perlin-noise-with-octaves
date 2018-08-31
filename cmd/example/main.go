package main

import (
	"fmt"
	"math"
)

var seed float64 = 100

type Perlin struct {
	X, Amp, Wavelength, Frequency, A, B float64
	Prng                                PRNG
	Pos                                 []float64
}

func main() {
	// demoPRNG()
	// demoInterpolate()

	amp := 128
	wl := 128
	octaves := 8
	divisor := 2
	w := 10
	noise := generateNoise(amp, wl, octaves, divisor, w)
	// for _, n := range noise {
	// 	fmt.Println(n.Pos)
	// }

	combined := combineNoise(noise)
	fmt.Println(combined)

}

func demoPRNG() {
	x := prng()
	fmt.Printf("Pseudo-random number generator example (seed:%f): %f\n", seed, x)
}

func demoInterpolate() {

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

func NewPerlin(amp, wl, width float64) Perlin {
	p := Perlin{
		X:          0,
		Amp:        amp,
		Wavelength: wl,
		Frequency:  1 / wl,
		Pos:        []float64{},
	}

	psng := PRNG{}
	p.A = psng.Next()
	p.B = psng.Next()

	for p.X < width {
		if math.Mod(p.X, p.Wavelength) == 0 {
			p.A = p.B
			p.B = psng.Next()
			p.Pos = append(p.Pos, p.A*p.Amp)
		} else {
			mu := (math.Mod(p.X, p.Wavelength) / p.Wavelength) * p.Amp
			interpolated := interpolate(p.A, p.B, mu)
			p.Pos = append(p.Pos, interpolated)
		}
		p.X++
	}
	return p
}

// PRNG is a psuedo-random number generator (linear congruential)
type PRNG struct {
	Z float64
}

func (p *PRNG) Next() float64 {
	p.Z = math.Mod(A*p.Z+C, M)
	return p.Z/M - 0.5
}

func generateNoise(amp, wl, octaves, divisor, width float64) []Perlin {
	result := []Perlin{}

	for i := 0.0; i < octaves; i++ {
		p := NewPerlin(amp, wl, width)
		amp = amp / divisor
		wl = wl / divisor
		result = append(result, p)
	}

	return result
}

/*
//combines octaves together
function CombineNoise(pl){
	var result = {pos: []};
	for(var i = 0, total = 0, j = 0; i < pl[0].pos.length; i++){
		total = 0;
		for(j = 0; j < pl.length; j++){
			total += pl[j].pos[i];
		}
		result.pos.push(total);
	}
	return result;
}
*/

func combineNoise(pl []Perlin) []float64 {
	combined := []float64{}

	for i, total, j := 0, 0.0, 0; i < len(pl[0].Pos); i++ {
		fmt.Println(i, total, j)
		total = 0
		for j = 0; j < len(pl); j++ {
			fmt.Println("\t", i, total, j)
			total += pl[j].Pos[i]
		}
		combined = append(combined, total)
	}

	return combined
}

/*

//perlin line plotting
function DrawLine(L){
	ctx.moveTo(0, h / 2);
	for(var i = 0; i < L.pos.length; i++){
		ctx.lineTo(i, h / 2 + L.pos[i]);
	}
	ctx.stroke();
}
*/
