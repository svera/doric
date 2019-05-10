package columns

import (
	"math/rand"
)

const maxTile = 6

type Piece struct {
	tiles []int
	Coords
	pit *Pit
}

func NewPiece(pit *Pit) *Piece {
	tiles := make([]int, 3)
	r := rand.New(source)

	tiles[0] = r.Intn(maxTile) + 1
	tiles[1] = r.Intn(maxTile) + 1
	tiles[2] = r.Intn(maxTile) + 1

	return &Piece{
		tiles,
		Coords{3, 0},
		pit,
	}
}

// X return the position of the piece in the X (horizontal) axis
func (p *Piece) X() int {
	return p.x
}

// Y return the position of the piece in the Y (vertical) axis
func (p *Piece) Y() int {
	return p.y
}

// Left move the piece to the left in the pit
func (p *Piece) Left() {
	if p.x > 0 {
		p.x--
	}
}

// Right move the piece to the right in the pit
func (p *Piece) Right() {
	if p.x < p.pit.Width()-1 {
		p.x++
	}
}

// Down move the current piece down in the pit. If the piece cannot fall further, returns false.
func (p *Piece) Down() bool {
	if p.y < p.pit.Height()-1 && p.pit.Cell(p.x, p.y+1) == Empty {
		p.y++
		return true
	}
	return false
}

// Rotate rotates piece tiles down. Last tile is moved to the first one
func (p *Piece) Rotate() {
	p.tiles[0], p.tiles[2] = p.tiles[2], p.tiles[0]
	p.tiles[1], p.tiles[2] = p.tiles[2], p.tiles[1]
}

func (p *Piece) Tiles() []int {
	return p.tiles
}

func (p *Piece) Renew() {
	p.tiles[0] = rand.Intn(maxTile) + 1
	p.tiles[1] = rand.Intn(maxTile) + 1
	p.tiles[2] = rand.Intn(maxTile) + 1
}

func (p *Piece) Copy(next *Piece) {
	p.tiles[0] = next.Tiles()[0]
	p.tiles[1] = next.Tiles()[1]
	p.tiles[2] = next.Tiles()[2]
	p.x = 3 // remove hardcoding
	p.y = 0
}
