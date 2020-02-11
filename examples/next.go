package main

import (
	"sync"

	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric"
)

// Next is an entity used to show next piece on screen
type Next struct {
	*tl.Entity
	Piece   *doric.Piece
	offsetX int
	offsetY int
	mux     sync.Locker
}

// NewNext returns a new Next instance
func NewNext(p *doric.Piece, offsetX, offsetY int, mux sync.Locker) *Next {
	return &Next{
		Entity:  tl.NewEntity(offsetX, offsetY, 1, 3),
		Piece:   p,
		offsetX: offsetX,
		offsetY: offsetY,
		mux:     mux,
	}
}

// Draw prints next piece on screen
func (n *Next) Draw(screen *tl.Screen) {
	n.mux.Lock()
	for i := range n.Piece.Tiles {
		screen.RenderCell(n.offsetX, n.offsetY-i, &tl.Cell{
			Bg: colors[n.Piece.Tiles[i]],
			Fg: tl.ColorBlack,
			Ch: '[',
		})
		screen.RenderCell(n.offsetX+1, n.offsetY-i, &tl.Cell{
			Bg: colors[n.Piece.Tiles[i]],
			Fg: tl.ColorBlack,
			Ch: ']',
		})
	}
	n.mux.Unlock()
}