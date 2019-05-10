package main

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric/pkg/columns"
)

type Pit struct {
	*tl.Entity
	pit     *columns.Pit
	offsetX int
	offsetY int
}

func NewPit(p *columns.Pit, offsetX int, offsetY int) *Pit {
	return &Pit{
		Entity:  tl.NewEntity(offsetX, offsetY, p.Width(), p.Height()),
		pit:     p,
		offsetX: offsetX,
		offsetY: offsetY,
	}
}

func (p *Pit) Draw(screen *tl.Screen) {
	// Pit bottom corners
	screen.RenderCell(p.offsetX, p.offsetY+p.pit.Height(), &tl.Cell{
		Bg: tl.ColorWhite,
		Ch: ' ',
	})
	screen.RenderCell(p.offsetX+p.pit.Width()+1, p.offsetY+p.pit.Height(), &tl.Cell{
		Bg: tl.ColorWhite,
		Ch: ' ',
	})

	for y := 0; y < p.pit.Height(); y++ {
		for x := 0; x < p.pit.Width(); x++ {
			// Pit left border
			if x == 0 {
				screen.RenderCell(p.offsetX, p.offsetY+y, &tl.Cell{
					Bg: tl.ColorWhite,
					Ch: ' ',
				})
			}
			// Pit right border
			if x == p.pit.Width()-1 {
				screen.RenderCell(p.offsetX+p.pit.Width()+1, p.offsetY+y, &tl.Cell{
					Bg: tl.ColorWhite,
					Ch: ' ',
				})
			}
			// Pit bottom
			if y == p.pit.Height()-1 {
				screen.RenderCell(p.offsetX+x+1, p.offsetY+y+1, &tl.Cell{
					Bg: tl.ColorWhite,
					Ch: ' ',
				})
			}
			if p.pit.Cell(x, y) == columns.Empty {
				screen.RenderCell(p.offsetX+x+1, p.offsetY+y, &tl.Cell{
					Bg: colors[p.pit.Cell(x, y)],
					Ch: ' ',
				})
			} else {
				screen.RenderCell(p.offsetX+x+1, p.offsetY+y, &tl.Cell{
					Bg: colors[p.pit.Cell(x, y)],
					Fg: tl.ColorWhite,
					Ch: ' ',
				})
			}
		}
	}
}
