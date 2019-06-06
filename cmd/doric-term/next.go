package main

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric/pkg/columns"
)

// Next is an entity used to show next piece on screen
type Next struct {
	*tl.Entity
	piece   *columns.Piece
	offsetX int
	offsetY int
}

// NewNext returns a new Next instance
func NewNext(p *columns.Piece, offsetX int, offsetY int) *Next {
	return &Next{
		Entity:  tl.NewEntity(offsetX, offsetY, 1, 3),
		piece:   p,
		offsetX: offsetX,
		offsetY: offsetY,
	}
}

// Draw prints next piece on screen
func (n *Next) Draw(screen *tl.Screen) {
	for i := range n.piece.Tiles() {
		screen.RenderCell(n.offsetX, n.offsetY-i, &tl.Cell{
			Bg: colors[n.piece.Tiles()[i]],
			Fg: tl.ColorBlack,
			Ch: chars[n.piece.Tiles()[i]],
		})
	}
}
