package columns

// values that represent empty or removable tiles in the pit
const (
	Remove = -1
	Empty  = 0
)

// Pit represents the field of play, holding the tiles that are falling
type Pit struct {
	width  int
	height int
	cells  [][]int
}

// NewPit return a new Pit instance
func NewPit(rows, cols int) *Pit {
	p := Pit{
		width:  cols,
		height: rows,
		cells:  make([][]int, rows),
	}
	for i := range p.cells {
		p.cells[i] = make([]int, cols)
	}
	return &p
}

// Reset empties the pit
func (p *Pit) Reset() {
	for y, row := range p.cells {
		for x := range row {
			p.cells[y][x] = Empty
		}
	}
}

// CheckLines scans pit lines looking for tiles to be removed.
// Tiles repeated in 3 or more consecutive positions horizontally, vertically or diagonally are to be removed.
func (p *Pit) CheckLines() int {
	remove := map[Coords]struct{}{}
	p.checkHorizontalLines(remove)
	p.checkVerticalLines(remove)
	p.checkDiagonalLines(remove)
	for coords := range remove {
		// Cells with negative values are cells with tiles to be removed
		p.cells[coords.y][coords.x] = Remove
	}
	return len(remove)
}

func (p *Pit) checkHorizontalLines(remove map[Coords]struct{}) {
	for y := p.height - 1; y >= 0; y-- {
		for x := 0; x < p.width-2; x++ {
			if p.cells[y][x] == Empty {
				continue
			}
			if p.cells[y][x] == p.cells[y][x+1] && p.cells[y][x+1] == p.cells[y][x+2] {
				remove[Coords{x, y}] = struct{}{}
				remove[Coords{x + 1, y}] = struct{}{}
				remove[Coords{x + 2, y}] = struct{}{}
			}
		}
	}
}

func (p *Pit) checkVerticalLines(remove map[Coords]struct{}) {
	for x := 0; x < p.width; x++ {
		for y := p.height - 1; y > 1; y-- {
			if p.cells[y][x] == Empty {
				break
			}
			if p.cells[y][x] == p.cells[y-1][x] && p.cells[y-1][x] == p.cells[y-2][x] {
				remove[Coords{x, y}] = struct{}{}
				remove[Coords{x, y - 1}] = struct{}{}
				remove[Coords{x, y - 2}] = struct{}{}
			}
		}
	}
}

func (p *Pit) checkDiagonalLines(remove map[Coords]struct{}) {
	for y := p.height - 1; y > 1; y-- {
		// Checks for tiles to be removed in diagonal / lines
		for x := 0; x < p.width-2 && y > 1; x++ {
			if p.cells[y][x] == Empty {
				continue
			}
			if p.cells[y][x] == p.cells[y-1][x+1] && p.cells[y-1][x+1] == p.cells[y-2][x+2] {
				remove[Coords{x, y}] = struct{}{}
				remove[Coords{x + 1, y - 1}] = struct{}{}
				remove[Coords{x + 2, y - 2}] = struct{}{}
			}
		}
		// Checks for tiles to be removed in diagonal \ lines
		for x := p.width - 1; x > 1 && y > 1; x-- {
			if p.cells[y][x] == Empty {
				continue
			}
			if p.cells[y][x] == p.cells[y-1][x-1] && p.cells[y-1][x-1] == p.cells[y-2][x-2] {
				remove[Coords{x, y}] = struct{}{}
				remove[Coords{x - 1, y - 1}] = struct{}{}
				remove[Coords{x - 2, y - 2}] = struct{}{}
			}
		}
	}
}

// Cell returns the passed coordinates cell value
func (p *Pit) Cell(x, y int) int {
	return p.cells[y][x]
}

// Width returns pit's width
func (p *Pit) Width() int {
	return p.width
}

// Height returns pit's height
func (p *Pit) Height() int {
	return p.height
}

// Settle moves down all tiles which have empty cells below
func (p *Pit) Settle() {
	for x := 0; x < p.Width(); x++ {
		tiles := []int{}
		for y := p.Height() - 1; y >= 0; y-- {
			// This cell contains a tile to be removed, do not put it in the slice of tiles to settle
			if p.cells[y][x] < 0 {
				continue
			}
			// There are no more tiles over an empty cell, so we can settle this column
			if p.cells[y][x] == Empty {
				for i := 0; i < p.Height(); i++ {
					if len(tiles)-1 >= i {
						p.cells[p.Height()-1-i][x] = tiles[i]
					} else {
						p.cells[p.Height()-1-i][x] = Empty
					}
				}
				break
			}
			tiles = append(tiles, p.cells[y][x])
		}
	}
}

// Consolidate put the values of the passed piece in the pit
func (p *Pit) Consolidate(pc *Piece) {
	for i, tile := range pc.Tiles() {
		if pc.Y()-i < 0 {
			return
		}
		p.cells[pc.Y()-i][pc.X()] = tile
	}
}
