package main

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric"
)

// Pit represents a pit on screen following Termloop's Drawable interface
type Pit struct {
	*tl.Entity
	Pit     doric.Pit
	width   int
	height  int
	offsetX int
	offsetY int
}

// NewPit returns a new pit instance
func NewPit(offsetX, offsetY, width, height int) *Pit {
	return &Pit{
		Entity:  tl.NewEntity(offsetX, offsetY, width, height),
		width:   width,
		height:  height,
		offsetX: offsetX,
		offsetY: offsetY,
	}
}

// Draw draws pit on screen
func (p *Pit) Draw(screen *tl.Screen) {
	var x, y int

	for y = 0; y <= p.height; y++ {
		for x = 0; x <= p.width; x++ {
			// Pit left border
			if x == 0 {
				screen.RenderCell(p.offsetX, p.offsetY+y, &tl.Cell{
					Bg: tl.ColorWhite,
					Ch: ' ',
				})
			}
			// Pit right border
			if x == p.width {
				screen.RenderCell(p.offsetX+p.width*2+1, p.offsetY+y, &tl.Cell{
					Bg: tl.ColorWhite,
					Ch: ' ',
				})
				continue
			}
			// Pit bottom
			if y == p.height {
				screen.RenderCell(p.offsetX+(x*2)+1, p.offsetY+y, &tl.Cell{
					Bg: tl.ColorWhite,
					Ch: ' ',
				})
				screen.RenderCell(p.offsetX+(x*2)+2, p.offsetY+y, &tl.Cell{
					Bg: tl.ColorWhite,
					Ch: ' ',
				})
				continue
			}
			// Tiles
			if p.Pit[y][x] >= doric.Empty {
				p.renderTile(screen, x, y)
			}
		}
	}
}

func (p *Pit) renderTile(screen *tl.Screen, x, y int) {
	leftCh, rightCh := '[', ']'

	if p.Pit[y][x] == doric.Empty {
		leftCh, rightCh = ' ', ' '
	}
	screen.RenderCell(p.offsetX+(x*2)+1, p.offsetY+y, &tl.Cell{
		Bg: colors[p.Pit[y][x]],
		Fg: tl.ColorBlack,
		Ch: leftCh,
	})
	screen.RenderCell(p.offsetX+(x*2)+2, p.offsetY+y, &tl.Cell{
		Bg: colors[p.Pit[y][x]],
		Fg: tl.ColorBlack,
		Ch: rightCh,
	})
}
