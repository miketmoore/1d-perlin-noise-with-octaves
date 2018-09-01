package main

import (
	"fmt"
	"math"
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
var w float64 = 10

func foo() {
	noise := generateNoise(amp, wl, octaves, divisor, w)
	combined := combineNoise(noise)
	fmt.Println(combined)
}

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

	noise := generateNoise(amp, wl, octaves, divisor, w)
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
			imd.Push(pixel.V(0, height/2))
			for i, c := range combined {
				imd.Color = colornames.Black
				imd.EndShape = imdraw.SharpEndShape
				x := float64(i)
				y := float64(height/2) + c
				imd.Push(pixel.V(x, y))

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
