package main

import (
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	w float64 = 800
	h float64 = 500
	M         = 4294967296
	A         = 1664525
	C         = 1
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "1D Perlin Noise",
		Bounds: pixel.R(0, 0, w, h),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	state := "generate"

	var data []float64

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyQ) {
			os.Exit(1)
		}

		if win.JustPressed(pixelgl.KeySpace) {
			state = "generate"
		}

		if state == "generate" {
			rand.Seed(time.Now().UTC().UnixNano())
			dataNoise := GenerateNoise(128, 128, 8, 2, w)
			data = CombineNoise(dataNoise)
			state = "draw"
		}

		if state == "draw" {

			win.Clear(colornames.White)
			drawLine(win, data)
			state = "nothing"
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

func drawLine(win *pixelgl.Window, combined []float64) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Black
	y := h / 2.0
	imd.Push(pixel.V(0, y))
	for i := 0; i < len(combined); i++ {
		x := float64(i)
		y2 := y + combined[i]
		imd.Push(pixel.V(x, y2))
	}
	imd.Line(1)
	imd.Draw(win)
}

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
