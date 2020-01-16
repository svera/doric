package main

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric/pkg/columns"
)

// Pit represents a pit on screen following Termloop's Drawable interface
type Pit struct {
	*tl.Entity
	Pit     columns.Pit
	offsetX int
	offsetY int
}

// NewPit returns a new pit instance
func NewPit(p columns.Pit, offsetX int, offsetY int) *Pit {
	return &Pit{
		Pit:     p,
		Entity:  tl.NewEntity(offsetX, offsetY, p.Width(), p.Height()),
		offsetX: offsetX,
		offsetY: offsetY,
	}
}

// Draw draws pit on screen
func (p *Pit) Draw(screen *tl.Screen) {
	var x, y int

	for y = 0; y <= p.Pit.Height(); y++ {
		for x = 0; x <= p.Pit.Width(); x++ {
			// Pit left border
			if x == 0 {
				screen.RenderCell(p.offsetX, p.offsetY+y, &tl.Cell{
					Bg: tl.ColorWhite,
					Ch: ' ',
				})
			}
			// Pit right border
			if x == p.Pit.Width() {
				screen.RenderCell(p.offsetX+p.Pit.Width()*2+1, p.offsetY+y, &tl.Cell{
					Bg: tl.ColorWhite,
					Ch: ' ',
				})
				continue
			}
			// Pit bottom
			if y == p.Pit.Height() {
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
			if p.Pit.Cell(x, y) > columns.Empty {
				screen.RenderCell(p.offsetX+(x*2)+1, p.offsetY+y, &tl.Cell{
					Bg: colors[p.Pit.Cell(x, y)],
					Fg: tl.ColorBlack,
					Ch: '[',
				})
				screen.RenderCell(p.offsetX+(x*2)+2, p.offsetY+y, &tl.Cell{
					Bg: colors[p.Pit.Cell(x, y)],
					Fg: tl.ColorBlack,
					Ch: ']',
				})
			} else {
				screen.RenderCell(p.offsetX+(x*2)+1, p.offsetY+y, &tl.Cell{
					Bg: colors[p.Pit.Cell(x, y)], // Failed with index out of range, check
					Ch: ' ',
				})
				screen.RenderCell(p.offsetX+(x*2)+2, p.offsetY+y, &tl.Cell{
					Bg: colors[p.Pit.Cell(x, y)],
					Ch: ' ',
				})
			}
		}
	}
}
