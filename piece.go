package doric

// maxTile is the maximum tile value a piece can contain
const maxTile = 6

// Randomizer defines a required method to get random integer values.
// Each value will map to one of the n possible tile types, where n = 6 by default.
type Randomizer interface {
	Intn(n int) int
}

// Piece represents a piece to fall in the well
type Piece struct {
	// Tiles composing the piece. Tile at index 0 corresponds to upper one,
	// while tile at index 2 refers to the bottom one.
	Tiles [3]int
	// Position of the piece in the well, using its bottom tile as reference.
	X, Y int
}

// randomize assigns the piece three new tiles
func (p *Piece) randomize(r Randomizer) {
	p.Tiles[0] = r.Intn(maxTile) + 1
	p.Tiles[1] = r.Intn(maxTile) + 1
	p.Tiles[2] = r.Intn(maxTile) + 1
}

// left moves the piece to the left in the well if that position is empty
// and not out of bounds
func (p *Piece) left(well Well) {
	if p.X > 0 && well[p.X-1][p.Y] == Empty {
		p.X--
	}
}

// right moves the piece to the right in the well if that position is empty
// and not out of bounds
func (p *Piece) right(well Well) {
	if p.X < well.width()-1 && well[p.X+1][p.Y] == Empty {
		p.X++
	}
}

// down moves the current piece down in the well. If the piece cannot fall further, returns false.
func (p *Piece) down(well Well) bool {
	if p.Y < well.height()-1 && well[p.X][p.Y+1] == Empty {
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
