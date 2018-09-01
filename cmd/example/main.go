package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	noise "github.com/miketmoore/perlin-noise"
	"golang.org/x/image/colornames"
)

var (
	win     *pixelgl.Window
	amp     float64 = 128
	wl      float64 = 128
	octaves float64 = 8
	divisor float64 = 2
)

const (
	height = 500
	width  = 500
)

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
			data = generateData()
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

func generateData() []float64 {
	generatedNoise := noise.GenerateNoise(amp, wl, octaves, divisor, width)
	fmt.Println("total noise: ", len(generatedNoise))
	for i, n := range generatedNoise {
		fmt.Println(i, len(n.Pos))
	}
	return noise.CombineNoise(generatedNoise)
}

func drawLine(win *pixelgl.Window, combined []float64) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Black
	y := float64(height / 2)
	imd.Push(pixel.V(0, y))
	for i, c := range combined {
		x := float64(i)
		y = y + c
		fmt.Println(i)
		imd.Push(pixel.V(x, y))
	}
	imd.Line(1)
	imd.Draw(win)
}
