package main

import (
	"sync"

	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric"
)

// Well represents a well on screen implementing Termloop's Drawable interface
type Well struct {
	*tl.Entity
	Well    doric.Well
	width   int
	height  int
	offsetX int
	offsetY int
	mux     sync.Locker
}

// NewWell returns a new well instance
func NewWell(p doric.Well, offsetX, offsetY, height, width int, mux sync.Locker) *Well {
	return &Well{
		Entity:  tl.NewEntity(offsetX, offsetY, width, height),
		Well:    p,
		width:   width,
		height:  height,
		offsetX: offsetX,
		offsetY: offsetY,
		mux:     mux,
	}
}

// Draw draws well on screen
func (p *Well) Draw(screen *tl.Screen) {
	var x, y int

	// Top left corner
	screen.RenderCell(p.offsetX-1, p.offsetY-1, &tl.Cell{
		Bg: tl.ColorWhite,
		Ch: ' ',
	})
	screen.RenderCell(p.offsetX, p.offsetY-1, &tl.Cell{
		Bg: tl.ColorWhite,
		Ch: ' ',
	})

	// Top right corner
	screen.RenderCell(p.offsetX+p.width*2+1, p.offsetY-1, &tl.Cell{
		Bg: tl.ColorWhite,
		Ch: ' ',
	})
	screen.RenderCell(p.offsetX+p.width*2+2, p.offsetY-1, &tl.Cell{
		Bg: tl.ColorWhite,
		Ch: ' ',
	})

	for y = 0; y <= p.height; y++ {
		for x = 0; x <= p.width; x++ {
			// Well left border
			if x == 0 {
				screen.RenderCell(p.offsetX-1, p.offsetY+y, &tl.Cell{
					Bg: tl.ColorWhite,
					Ch: ' ',
				})
				screen.RenderCell(p.offsetX, p.offsetY+y, &tl.Cell{
					Bg: tl.ColorWhite,
					Ch: ' ',
				})
			}
			// Well right border
			if x == p.width {
				screen.RenderCell(p.offsetX+p.width*2+1, p.offsetY+y, &tl.Cell{
					Bg: tl.ColorWhite,
					Ch: ' ',
				})
				screen.RenderCell(p.offsetX+p.width*2+2, p.offsetY+y, &tl.Cell{
					Bg: tl.ColorWhite,
					Ch: ' ',
				})
				continue
			}
			// Well top
			if y == 0 {
				screen.RenderCell(p.offsetX+(x*2)+1, p.offsetY+y-1, &tl.Cell{
					Bg: tl.ColorWhite,
					Ch: ' ',
				})
				screen.RenderCell(p.offsetX+(x*2)+2, p.offsetY+y-1, &tl.Cell{
					Bg: tl.ColorWhite,
					Ch: ' ',
				})
				continue
			}
			// Well bottom
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
			if p.Well[x][y] >= doric.Empty {
				p.renderTile(screen, x, y)
			}
			p.mux.Unlock()
		}
	}
}

func (p *Well) renderTile(screen *tl.Screen, x, y int) {
	leftCh, rightCh := '[', ']'

	if p.Well[x][y] == doric.Empty {
		leftCh, rightCh = ' ', ' '
	}
	screen.RenderCell(p.offsetX+(x*2)+1, p.offsetY+y, &tl.Cell{
		Bg: colors[p.Well[x][y]],
		Fg: tl.ColorBlack,
		Ch: leftCh,
	})
	screen.RenderCell(p.offsetX+(x*2)+2, p.offsetY+y, &tl.Cell{
		Bg: colors[p.Well[x][y]],
		Fg: tl.ColorBlack,
		Ch: rightCh,
	})
}
