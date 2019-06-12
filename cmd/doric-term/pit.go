package main

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric/pkg/columns"
)

// Pit represents a pit on screen following Termloop's Drawable interface
type Pit struct {
	*tl.Entity
	pit     *columns.Pit
	offsetX int
	offsetY int
}

// NewPit returns a new pit instance
func NewPit(p *columns.Pit, offsetX int, offsetY int) *Pit {
	return &Pit{
		Entity:  tl.NewEntity(offsetX, offsetY, p.Width(), p.Height()),
		pit:     p,
		offsetX: offsetX,
		offsetY: offsetY,
	}
}

// Draw draws pit on screen
func (p *Pit) Draw(screen *tl.Screen) {
	var x, y int

	for y = 0; y <= p.pit.Height(); y++ {
		for x = 0; x <= p.pit.Width(); x++ {
			// Pit left border
			if x == 0 {
				screen.RenderCell(p.offsetX, p.offsetY+y, &tl.Cell{
					Bg: tl.ColorWhite,
					Ch: ' ',
				})
			}
			// Pit right border
			if x == p.pit.Width() {
				screen.RenderCell(p.offsetX+p.pit.Width()*2+1, p.offsetY+y, &tl.Cell{
					Bg: tl.ColorWhite,
					Ch: ' ',
				})
				continue
			}
			// Pit bottom
			if y == p.pit.Height() {
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
			if p.pit.Cell(x, y) > columns.Empty {
				screen.RenderCell(p.offsetX+(x*2)+1, p.offsetY+y, &tl.Cell{
					Bg: colors[p.pit.Cell(x, y)],
					Fg: tl.ColorBlack,
					Ch: '[',
				})
				screen.RenderCell(p.offsetX+(x*2)+2, p.offsetY+y, &tl.Cell{
					Bg: colors[p.pit.Cell(x, y)],
					Fg: tl.ColorBlack,
					Ch: ']',
				})
			} else {
				screen.RenderCell(p.offsetX+(x*2)+1, p.offsetY+y, &tl.Cell{
					Bg: colors[p.pit.Cell(x, y)],
					Ch: ' ',
				})
				screen.RenderCell(p.offsetX+(x*2)+2, p.offsetY+y, &tl.Cell{
					Bg: colors[p.pit.Cell(x, y)],
					Ch: ' ',
				})
			}
		}
	}
}
