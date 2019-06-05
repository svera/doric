package columns

import (
	"math/rand"
)

// maxTile is the maximum tile value a piece can contain
const maxTile = 6

// Piece represents a piece to fall in the pit
type Piece struct {
	tiles []int
	Coords
	pit *Pit
}

// NewPiece returns a new Piece instance
func NewPiece(pit *Pit) *Piece {
	return &Piece{
		make([]int, 3),
		Coords{3, 0},
		pit,
	}
}

// Reset assigns three new tiles to the piece, and resets its position to the initial one
func (p *Piece) Reset() {
	p.randomizeTiles()
	p.x = p.pit.Width() / 2
	p.y = 0
}

func (p *Piece) randomizeTiles() {
	r := rand.New(source)

	p.tiles[0] = r.Intn(maxTile) + 1
	p.tiles[1] = r.Intn(maxTile) + 1
	p.tiles[2] = r.Intn(maxTile) + 1
}

// X returns the position of the piece in the X (horizontal) axis
func (p *Piece) X() int {
	return p.x
}

// Y returns the position of the piece in the Y (vertical) axis
func (p *Piece) Y() int {
	return p.y
}

// Left moves the piece to the left in the pit
func (p *Piece) Left() {
	if p.x > 0 {
		p.x--
	}
}

// Right moves the piece to the right in the pit
func (p *Piece) Right() {
	if p.x < p.pit.Width()-1 {
		p.x++
	}
}

// Down moves the current piece down in the pit. If the piece cannot fall further, returns false.
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

// Tiles returns Piece's tiles
func (p *Piece) Tiles() []int {
	return p.tiles
}

// Renew assigns the piece three new tiles
func (p *Piece) Renew() {
	p.tiles[0] = rand.Intn(maxTile) + 1
	p.tiles[1] = rand.Intn(maxTile) + 1
	p.tiles[2] = rand.Intn(maxTile) + 1
}

// Copy copies the tiles from the passed piece, and resets its position to the initial one
func (p *Piece) Copy(next *Piece) {
	p.tiles[0] = next.Tiles()[0]
	p.tiles[1] = next.Tiles()[1]
	p.tiles[2] = next.Tiles()[2]
	p.x = p.pit.Width() / 2
	p.y = 0
}
