package main

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/kbinani/screenshot"

	"github.com/supermuesli/anim8/pkg/render"
)

func run() {
	// get display dimensions
	bounds := screenshot.GetDisplayBounds(0)
	width := float64(bounds.Dx())*0.9
	height := float64(bounds.Dy())*0.9

	// get bindata stuff
	brush, err := Asset("data/brush.png")
	if err != nil {
		panic(err)
	}

	font, err := Asset("data/ka1.ttf")
	if err != nil {
		panic(err)
	}

	// initialize new canvas
	canvas := render.NewCanvas(width, height, brush, font)

	// render loop
	for !canvas.Win.Closed() {
		canvas.Poll()
		canvas.Draw()		
		
		// manually enforce FPS
		<-canvas.FPS
	}
}

func main() {
	pixelgl.Run(run)
}