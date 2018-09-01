package noise

import (
	"math"
	"math/rand"
)

type Perlin struct {
	X, Amp, Wavelength, Frequency, A, B float64
	Prng                                PRNG
	Pos                                 []float64
}

/*
function Interpolate(pa, pb, px){
	var ft = px * Math.PI,
		f = (1 - Math.cos(ft)) * 0.5;
	return pa * (1 - f) + pb * f;
}
*/
func interpolate(pa, pb, px float64) float64 {
	ft := px * math.Pi
	f := (1 - math.Cos(ft)) * 0.5
	return pa*(1-f) + pb*f
}

/*
//1D perlin line generator
function Perlin(amp, wl, width){
	this.x = 0;
	this.amp = amp;
	this.wl = wl;
	this.fq = 1 / wl;
	this.psng = new PSNG();
	this.a = this.psng.next();
	this.b = this.psng.next();
	this.pos = [];
	while(this.x < width){
		if(this.x % this.wl === 0){
			this.a = this.b;
			this.b = this.psng.next();
			this.pos.push(this.a * this.amp);
		}else{
			this.pos.push(Interpolate(this.a, this.b, (this.x % this.wl) / this.wl) * this.amp);
		}
		this.x++;
	}
}
*/
func NewPerlin(amp, wl, width float64) Perlin {
	r := NewPRNG()
	p := Perlin{
		X:          0,
		Amp:        amp,
		Wavelength: wl,
		Frequency:  1 / wl,
		Prng:       r,
		A:          r.Next(),
		B:          r.Next(),
		Pos:        []float64{},
	}
	for p.X < width {
		if math.Mod(p.X, p.Wavelength) == 0 {
			p.A = p.B
			p.B = p.Prng.Next()
			p.Pos = append(p.Pos, p.A*p.Amp)
		} else {
			p.Pos = append(p.Pos, interpolate(p.A, p.B, (math.Mod(p.X, p.Wavelength)/p.Wavelength)*p.Amp))
		}
		p.X++
	}
	return p
}

var seed float64 = 100
var M float64 = 4294967296

// a - 1 should be divisible by m's prime factors
var A float64 = 1664525

// c and m should be co-prime
var C float64 = 1

// var Z float64 = math.Floor(seed)

// PRNG is a psuedo-random number generator (linear congruential)
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

/*
//octave generator
function GenerateNoise(amp, wl, octaves, divisor, width){
	var result = [];
	for(var i = 0; i < octaves; i++){
		result.push(new Perlin(amp, wl, width));
		amp /= divisor;
		wl /= divisor;
	}
	return result;
}
*/
func GenerateNoise(amp, wl, octaves, divisor, width float64) []Perlin {
	result := []Perlin{}

	for i := 0; i < int(octaves); i++ {
		p := NewPerlin(amp, wl, width)
		result = append(result, p)
		amp = amp / divisor
		wl = wl / divisor
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
func CombineNoise(pl []Perlin) []float64 {
	combined := []float64{}

	for i, total, j := 0, 0.0, 0; i < len(pl[0].Pos); i++ {
		total = 0
		for j = 0; j < len(pl); j++ {
			total += pl[j].Pos[i]
		}
		combined = append(combined, total)
	}

	return combined
}
