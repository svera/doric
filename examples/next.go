package main

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric"
)

// Next is an entity used to show next piece on screen
type Next struct {
	*tl.Entity
	Piece   *doric.Piece
	offsetX int
	offsetY int
}

// NewNext returns a new Next instance
func NewNext(offsetX int, offsetY int) *Next {
	return &Next{
		Entity:  tl.NewEntity(offsetX, offsetY, 1, 3),
		offsetX: offsetX,
		offsetY: offsetY,
	}
}

// Draw prints next piece on screen
func (n *Next) Draw(screen *tl.Screen) {
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
}
