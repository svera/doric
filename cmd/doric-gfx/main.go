package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Columns",
		Bounds: pixel.R(0, 0, 960, 672),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	pic, err := loadPicture("../../assets/backgrounds.png")
	if err != nil {
		panic(err)
	}
	jewels, err := loadPicture("../../assets/jewels.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pixel.R(326, 229, 646, 453))
	orange := pixel.NewSprite(jewels, pixel.R(19, 151, 33, 167))
	mat := pixel.IM
	mat = mat.Moved(win.Bounds().Center())
	mat = mat.ScaledXY(win.Bounds().Center(), pixel.V(3, 3))

	sprite.Draw(win, mat)
	orange.Draw(win, mat)

	for !win.Closed() {
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
