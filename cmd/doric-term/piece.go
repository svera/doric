package main

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric/pkg/columns"
)

type Piece struct {
	*tl.Entity
	piece   *columns.Piece
	offsetX int
	offsetY int
}

func NewPiece(p *columns.Piece, offsetX int, offsetY int) *Piece {
	return &Piece{
		Entity:  tl.NewEntity(p.X()+offsetX+1, p.Y()+offsetY, 1, 3),
		piece:   p,
		offsetX: offsetX,
		offsetY: offsetY,
	}
}

func (p *Piece) Draw(screen *tl.Screen) {
	for i := range p.piece.Tiles() {
		if i > p.piece.Y() {
			continue
		}
		screen.RenderCell(p.piece.X()+p.offsetX+1, p.piece.Y()+p.offsetY-i, &tl.Cell{
			Bg: colors[p.piece.Tiles()[i]],
			Fg: tl.ColorBlack,
			Ch: '·',
		})
	}
}

func (p *Piece) Tick(event tl.Event) {
	if event.Type == tl.EventKey { // Is it a keyboard event?
		switch event.Key { // If so, switch on the pressed key.
		case tl.KeyArrowRight:
			p.piece.Right()
		case tl.KeyArrowLeft:
			p.piece.Left()
		case tl.KeyArrowDown:
			p.piece.Down()
		case tl.KeyTab:
			p.piece.Rotate()
		}
	}
}
