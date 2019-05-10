package main

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric/pkg/columns"
)

type Next struct {
	*tl.Entity
	piece   *columns.Piece
	offsetX int
	offsetY int
}

func NewNext(p *columns.Piece, offsetX int, offsetY int) *Next {
	return &Next{
		Entity:  tl.NewEntity(offsetX, offsetY, 1, 3),
		piece:   p,
		offsetX: offsetX,
		offsetY: offsetY,
	}
}

func (n *Next) Draw(screen *tl.Screen) {
	for i := range n.piece.Tiles() {
		screen.RenderCell(n.offsetX, n.offsetY-i, &tl.Cell{
			Bg: colors[n.piece.Tiles()[i]],
			Ch: ' ',
		})
	}
}
