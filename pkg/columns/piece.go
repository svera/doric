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
		Coords{pit.width / 2, 0},
		pit,
	}
}

// Reset assigns three new tiles to the piece, and resets its position to the initial one
func (p *Piece) Reset() {
	p.RandomizeTiles()
	p.x = p.pit.width / 2
	p.y = 0
}

// RandomizeTiles assigns the piece three new tiles
func (p *Piece) RandomizeTiles() {
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

// Left moves the piece to the left in the pit if that position is empty
// and not out of bounds
func (p *Piece) Left() {
	if p.x > 0 && p.pit.Cell(p.x-1, p.y) == Empty {
		p.x--
	}
}

// Right moves the piece to the right in the pit if that position is empty
// and not out of bounds
func (p *Piece) Right() {
	if p.x < p.pit.width-1 && p.pit.Cell(p.x+1, p.y) == Empty {
		p.x++
	}
}

// Down moves the current piece down in the pit. If the piece cannot fall further, returns false.
func (p *Piece) Down() bool {
	if p.y < p.pit.height-1 && p.pit.Cell(p.x, p.y+1) == Empty {
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

// Copy copies the tiles from the passed piece, and resets its position to the initial one
func (p *Piece) Copy(next *Piece) {
	p.tiles[0] = next.Tiles()[0]
	p.tiles[1] = next.Tiles()[1]
	p.tiles[2] = next.Tiles()[2]
	p.x = p.pit.width / 2
	p.y = 0
}
