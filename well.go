package doric

// Values that represent empty or removable tiles in the well
const (
	Remove = -1
	Empty  = 0
)

// Standard Well dimensions (in tiles) as per commercial SEGA versions
const (
	StandardWidth  = 6
	StandardHeight = 13
)

// coords represent the coordinates of a tile or cell in the well
type coords struct {
	x int
	y int
}

// Well is a slice of slices which represents the field of play, holding the tiles that are falling.
// First index represents tiles in the X (horizontal) axis, second index refers to the Y (vertical) axis.
type Well [][]int

// NewWell return a new empty Well instance
func NewWell(rows, cols int) Well {
	var p Well
	p = make([][]int, cols)
	for i := range p {
		p[i] = make([]int, rows)
	}
	return p
}

// markTilesToRemove scans well lines looking for tiles to be removed, amd mark those tiles.
// Tiles repeated in 3 or more consecutive positions horizontally, vertically or diagonally are to be removed.
func (p Well) markTilesToRemove() int {
	remove := map[coords]struct{}{}
	p.checkHorizontalLines(remove)
	p.checkVerticalLines(remove)
	p.checkDiagonalLines(remove)
	for coords := range remove {
		// Cells with negative values are cells with tiles to be removed
		p[coords.x][coords.y] = Remove
	}
	return len(remove)
}

func (p Well) checkHorizontalLines(remove map[coords]struct{}) {
	for y := p.height() - 1; y >= 0; y-- {
		for x := 0; x < p.width()-2; x++ {
			if p[x][y] == Empty || p[x][y] == Remove {
				continue
			}
			if p[x][y] == p[x+1][y] && p[x+1][y] == p[x+2][y] {
				remove[coords{x, y}] = struct{}{}
				remove[coords{x + 1, y}] = struct{}{}
				remove[coords{x + 2, y}] = struct{}{}
			}
		}
	}
}

func (p Well) checkVerticalLines(remove map[coords]struct{}) {
	for x := 0; x < p.width(); x++ {
		for y := p.height() - 1; y > 1; y-- {
			if p[x][y] == Empty || p[x][y] == Remove {
				break
			}
			if p[x][y] == p[x][y-1] && p[x][y-1] == p[x][y-2] {
				remove[coords{x, y}] = struct{}{}
				remove[coords{x, y - 1}] = struct{}{}
				remove[coords{x, y - 2}] = struct{}{}
			}
		}
	}
}

func (p Well) checkDiagonalLines(remove map[coords]struct{}) {
	for y := p.height() - 1; y > 1; y-- {
		// Checks for tiles to be removed in diagonal / lines
		for x := 0; x < p.width()-2 && y > 1; x++ {
			if p[x][y] == Empty || p[x][y] == Remove {
				continue
			}
			if p[x][y] == p[x+1][y-1] && p[x+1][y-1] == p[x+2][y-2] {
				remove[coords{x, y}] = struct{}{}
				remove[coords{x + 1, y - 1}] = struct{}{}
				remove[coords{x + 2, y - 2}] = struct{}{}
			}
		}
		// Checks for tiles to be removed in diagonal \ lines
		for x := p.width() - 1; x > 1 && y > 1; x-- {
			if p[x][y] == Empty {
				continue
			}
			if p[x][y] == p[x-1][y-1] && p[x-1][y-1] == p[x-2][y-2] {
				remove[coords{x, y}] = struct{}{}
				remove[coords{x - 1, y - 1}] = struct{}{}
				remove[coords{x - 2, y - 2}] = struct{}{}
			}
		}
	}
}

// Width returns well's width
func (p Well) width() int {
	return len(p)
}

// Height returns well's height
func (p Well) height() int {
	return len(p[0])
}

// settle moves down all tiles which have empty cells below
func (p Well) settle() {
	for x := 0; x < p.width(); x++ {
		moveDown := 0
		for y := p.height() - 1; y >= 0; y-- {
			// This cell contains a tile to be removed, do not put it in the slice of tiles to settle again
			if p[x][y] < 0 {
				p[x][y] = Empty
				moveDown++
				continue
			}
			// There are no more tiles over an empty cell, so we can stop processing this column
			if p[x][y] == Empty {
				break
			}
			if moveDown > 0 {
				p[x][y+moveDown] = p[x][y]
				p[x][y] = Empty
			}
		}
	}
}

// lock put the values of the passed piece in the well
func (p Well) lock(pc *Piece) {
	for i, tile := range pc.Tiles {
		if pc.Y-i < 0 {
			return
		}
		p[pc.X][pc.Y-i] = tile
	}
}

func (p Well) copy() Well {
	well := NewWell(p.height(), p.width())
	for i := range p {
		copy(well[i], p[i])
	}
	return well
}
