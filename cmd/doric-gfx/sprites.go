package main

import "github.com/faiface/pixel"

const (
	yellow = iota + 1
	orange
	green
	purple
	red
	blue
)

func loadJewelsSprites(sheet pixel.Picture) map[int]*pixel.Sprite {
	jewels := make(map[int]*pixel.Sprite, 6)

	jewels[yellow] = pixel.NewSprite(sheet, pixel.R(2, 152, 16, 167))
	jewels[orange] = pixel.NewSprite(sheet, pixel.R(19, 151, 33, 167))
	jewels[green] = pixel.NewSprite(sheet, pixel.R(22, 151, 37, 166))

	return jewels
}
