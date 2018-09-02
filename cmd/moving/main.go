package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	noise "github.com/miketmoore/perlin-noise"
	"golang.org/x/image/colornames"
)

const (
	w float64 = 800
	h float64 = 500
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

	dataA := generate()
	dataB := generate()

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyQ) {
			os.Exit(1)
		}

		if win.JustPressed(pixelgl.KeySpace) {
			state = "generate"
		}

		if state == "generate" {
			state = "draw"
		}

		if state == "draw" {
			win.Clear(colornames.White)
			drawLine(win, dataA)
			drawLine(win, dataB)
			state = "nothing"
		}

		win.Update()
	}
}

func generate() []float64 {
	rand.Seed(time.Now().UTC().UnixNano())
	dataNoise := noise.GenerateNoise(128, 128, 8, 2, w)
	return noise.CombineNoise(dataNoise)
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
	imd.Line(2)
	imd.Draw(win)
}
