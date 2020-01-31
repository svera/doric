package doric

// maxTile is the maximum tile value a piece can contain
const maxTile = 6

// Randomizer defines a required method to get random integer values
type Randomizer interface {
	Intn(n int) int
}

// Piece represents a piece to fall in the pit
type Piece struct {
	Tiles [3]int
	Coords
}

// NewPiece returns a new Piece instance
func NewPiece(r Randomizer) *Piece {
	p := &Piece{
		Tiles: [3]int{},
	}
	p.randomize(r)
	return p
}

// randomize assigns the piece three new tiles
func (p *Piece) randomize(r Randomizer) {
	p.Tiles[0] = r.Intn(maxTile) + 1
	p.Tiles[1] = r.Intn(maxTile) + 1
	p.Tiles[2] = r.Intn(maxTile) + 1
}

// left moves the piece to the left in the pit if that position is empty
// and not out of bounds
func (p *Piece) left(pit Pit) {
	if p.X > 0 && pit.Cell(p.X-1, p.Y) == Empty {
		p.X--
	}
}

// right moves the piece to the right in the pit if that position is empty
// and not out of bounds
func (p *Piece) right(pit Pit) {
	if p.X < pit.Width()-1 && pit.Cell(p.X+1, p.Y) == Empty {
		p.X++
	}
}

// down moves the current piece down in the pit. If the piece cannot fall further, returns false.
func (p *Piece) down(pit Pit) bool {
	if p.Y < pit.Height()-1 && pit.Cell(p.X, p.Y+1) == Empty {
		p.Y++
		return true
	}
	return false
}

// rotate rotates piece tiles down. Last tile is moved to the first one
func (p *Piece) rotate() {
	p.Tiles[0], p.Tiles[2] = p.Tiles[2], p.Tiles[0]
	p.Tiles[1], p.Tiles[2] = p.Tiles[2], p.Tiles[1]
}

// copy copies the tiles from the passed piece, and resets its position to the initial one
func (p *Piece) copy(next *Piece, col int) {
	p.Tiles[0] = next.Tiles[0]
	p.Tiles[1] = next.Tiles[1]
	p.Tiles[2] = next.Tiles[2]
	p.X = col
	p.Y = 0
}
