package main

import (
	"sync"

	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric"
)

// Pit represents a pit on screen implementing Termloop's Drawable interface
type Pit struct {
	*tl.Entity
	Pit     doric.Pit
	width   int
	height  int
	offsetX int
	offsetY int
	mux     sync.Locker
}

// NewPit returns a new pit instance
func NewPit(p doric.Pit, offsetX, offsetY, height, width int, mux sync.Locker) *Pit {
	return &Pit{
		Entity:  tl.NewEntity(offsetX, offsetY, width, height),
		Pit:     p,
		width:   width,
		height:  height,
		offsetX: offsetX,
		offsetY: offsetY,
		mux:     mux,
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
			p.mux.Lock()
			if p.Pit[x][y] >= doric.Empty {
				p.renderTile(screen, x, y)
			}
			p.mux.Unlock()
		}
	}
}

func (p *Pit) renderTile(screen *tl.Screen, x, y int) {
	leftCh, rightCh := '[', ']'

	if p.Pit[x][y] == doric.Empty {
		leftCh, rightCh = ' ', ' '
	}
	screen.RenderCell(p.offsetX+(x*2)+1, p.offsetY+y, &tl.Cell{
		Bg: colors[p.Pit[x][y]],
		Fg: tl.ColorBlack,
		Ch: leftCh,
	})
	screen.RenderCell(p.offsetX+(x*2)+2, p.offsetY+y, &tl.Cell{
		Bg: colors[p.Pit[x][y]],
		Fg: tl.ColorBlack,
		Ch: rightCh,
	})
}
