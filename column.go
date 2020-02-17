package doric

// maxTile is the maximum tile value a column can contain
const maxTile = 6

// TilesFactory defines the signature of the method to build a column tileset.
type TilesFactory func(int) [3]int

// Column represents a column to fall in the well
type Column struct {
	// Tiles composing the column. Tile at index 0 corresponds to upper one,
	// while tile at index 2 refers to the bottom one. Possible tile values go from
	// 1 to maxTile.
	Tiles [3]int
	// Position of the column in the well, using its bottom tile as reference.
	X, Y int
}

// left moves the column to the left in the well if that position is empty
// and not out of bounds
func (p *Column) left(well Well) {
	if p.X > 0 && well[p.X-1][p.Y] == Empty {
		p.X--
	}
}

// right moves the column to the right in the well if that position is empty
// and not out of bounds
func (p *Column) right(well Well) {
	if p.X < well.width()-1 && well[p.X+1][p.Y] == Empty {
		p.X++
	}
}

// down moves the current column down in the well. If the column cannot fall further, returns false.
func (p *Column) down(well Well) bool {
	if p.Y < well.height()-1 && well[p.X][p.Y+1] == Empty {
		p.Y++
		return true
	}
	return false
}

// rotate rotates column tiles down. Last tile is moved to the first one
func (p *Column) rotate() {
	p.Tiles[0], p.Tiles[2] = p.Tiles[2], p.Tiles[0]
	p.Tiles[1], p.Tiles[2] = p.Tiles[2], p.Tiles[1]
}

// copy copies the tiles from the passed column, and resets its position to the initial one
func (p *Column) copy(next *Column, col int) {
	p.Tiles[0] = next.Tiles[0]
	p.Tiles[1] = next.Tiles[1]
	p.Tiles[2] = next.Tiles[2]
	p.X = col
	p.Y = 0
}
