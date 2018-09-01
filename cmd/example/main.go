package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	seed float64 = 100
	win  *pixelgl.Window
)

type Perlin struct {
	X, Amp, Wavelength, Frequency, A, B float64
	Prng                                PRNG
	Pos                                 []float64
}

var amp float64 = 128
var wl float64 = 128
var octaves float64 = 8
var divisor float64 = 2

const height = 500
const width = 500

func run() {

	// Setup GUI window
	cfg := pixelgl.WindowConfig{
		Title:  "1D Perlin Noise",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	state := "draw"

	noise := generateNoise(amp, wl, octaves, divisor, width)
	fmt.Println("total noise: ", len(noise))
	// should be 1,168 values in n.Pos slice
	for i, n := range noise {
		fmt.Println(i, len(n.Pos))
	}
	combined := combineNoise(noise)

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyQ) {
			os.Exit(1)
		}

		if win.JustPressed(pixelgl.KeySpace) {
			state = "draw"
		}

		if state == "draw" {
			win.Clear(colornames.White)

			/*
				function DrawLine(L){
					ctx.moveTo(0, h / 2);
					for(var i = 0; i < L.pos.length; i++){
						ctx.lineTo(i, h / 2 + L.pos[i]);
					}
					ctx.stroke();
				}
			*/
			imd := imdraw.New(nil)
			imd.Color = colornames.Black
			imd.Push(pixel.V(0, height/2))
			for i, c := range combined {
				imd.EndShape = imdraw.SharpEndShape
				x := float64(i)
				y := float64(height/2) + (c * 10)
				imd.Push(pixel.V(x+50, y))

			}
			imd.Line(1)
			imd.Draw(win)
			state = "nothing"
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

func drawRect(win *pixelgl.Window, x1, y1, w, h float64, color pixel.RGBA) {
	rect := imdraw.New(nil)
	rect.Color = color
	rect.Push(pixel.V(x1, y1))
	rect.Push(pixel.V(x1+w, y1+h))
	rect.Rectangle(0)
	rect.Draw(win)
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
	r := PRNG{
		Z: math.Floor(rand.Float64() * M),
	}
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

var M float64 = 4294967296

// a - 1 should be divisible by m's prime factors
var A float64 = 1664525

// c and m should be co-prime
var C float64 = 1

var Z float64 = math.Floor(seed)

// PRNG is a psuedo-random number generator (linear congruential)
type PRNG struct {
	Z float64
}

/*
function PSNG(){
	this.Z = Math.floor(Math.random() * M);
	this.next = function(){
		this.Z = (A * this.Z + C) % M;
		return this.Z / M - 0.5;
	}
}
*/

func (p *PRNG) Next() float64 {
	p.Z = math.Mod(A*p.Z+C, M)
	return p.Z/M - 0.5
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
func generateNoise(amp, wl, octaves, divisor, width float64) []Perlin {
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
func combineNoise(pl []Perlin) []float64 {
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
