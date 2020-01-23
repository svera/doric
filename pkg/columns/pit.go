package columns

// values that represent empty or removable tiles in the pit
const (
	Remove = -1
	Empty  = 0
)

// Pit represents the field of play, holding the tiles that are falling
type Pit [][]int

// NewPit return a new Pit instance
func NewPit(rows, cols int) Pit {
	var p Pit
	p = make([][]int, rows)
	for i := range p {
		p[i] = make([]int, cols)
	}
	return p
}

// markTilesToRemove scans pit lines looking for tiles to be removed, amd mark those tiles.
// Tiles repeated in 3 or more consecutive positions horizontally, vertically or diagonally are to be removed.
func (p Pit) markTilesToRemove() int {
	remove := map[Coords]struct{}{}
	p.checkHorizontalLines(remove)
	p.checkVerticalLines(remove)
	p.checkDiagonalLines(remove)
	for coords := range remove {
		// Cells with negative values are cells with tiles to be removed
		p[coords.Y][coords.X] = Remove
	}
	return len(remove)
}

func (p Pit) checkHorizontalLines(remove map[Coords]struct{}) {
	for y := p.Height() - 1; y >= 0; y-- {
		for x := 0; x < p.Width()-2; x++ {
			if p[y][x] == Empty || p[y][x] == Remove {
				continue
			}
			if p[y][x] == p[y][x+1] && p[y][x+1] == p[y][x+2] {
				remove[Coords{x, y}] = struct{}{}
				remove[Coords{x + 1, y}] = struct{}{}
				remove[Coords{x + 2, y}] = struct{}{}
			}
		}
	}
}

func (p Pit) checkVerticalLines(remove map[Coords]struct{}) {
	for x := 0; x < p.Width(); x++ {
		for y := p.Height() - 1; y > 1; y-- {
			if p[y][x] == Empty || p[y][x] == Remove {
				break
			}
			if p[y][x] == p[y-1][x] && p[y-1][x] == p[y-2][x] {
				remove[Coords{x, y}] = struct{}{}
				remove[Coords{x, y - 1}] = struct{}{}
				remove[Coords{x, y - 2}] = struct{}{}
			}
		}
	}
}

func (p Pit) checkDiagonalLines(remove map[Coords]struct{}) {
	for y := p.Height() - 1; y > 1; y-- {
		// Checks for tiles to be removed in diagonal / lines
		for x := 0; x < p.Width()-2 && y > 1; x++ {
			if p[y][x] == Empty || p[y][x] == Remove {
				continue
			}
			if p[y][x] == p[y-1][x+1] && p[y-1][x+1] == p[y-2][x+2] {
				remove[Coords{x, y}] = struct{}{}
				remove[Coords{x + 1, y - 1}] = struct{}{}
				remove[Coords{x + 2, y - 2}] = struct{}{}
			}
		}
		// Checks for tiles to be removed in diagonal \ lines
		for x := p.Width() - 1; x > 1 && y > 1; x-- {
			if p[y][x] == Empty {
				continue
			}
			if p[y][x] == p[y-1][x-1] && p[y-1][x-1] == p[y-2][x-2] {
				remove[Coords{x, y}] = struct{}{}
				remove[Coords{x - 1, y - 1}] = struct{}{}
				remove[Coords{x - 2, y - 2}] = struct{}{}
			}
		}
	}
}

// Cell returns the passed coordinates cell value
func (p Pit) Cell(x, y int) int {
	return p[y][x]
}

// Width returns pit's width
func (p Pit) Width() int {
	return len(p[0])
}

// Height returns pit's height
func (p Pit) Height() int {
	return len(p)
}

// settle moves down all tiles which have empty cells below
func (p Pit) settle() {
	for x := 0; x < p.Width(); x++ {
		tiles := []int{}
		for y := p.Height() - 1; y >= 0; y-- {
			// This cell contains a tile to be removed, do not put it in the slice of tiles to settle again
			if p[y][x] < 0 {
				continue
			}
			// There are no more tiles over an empty cell, so we can settle this column
			if p[y][x] == Empty {
				for i := 0; i < p.Height(); i++ {
					if len(tiles)-1 >= i {
						p[p.Height()-1-i][x] = tiles[i]
					} else {
						p[p.Height()-1-i][x] = Empty
					}
				}
				break
			}
			tiles = append(tiles, p[y][x])
		}
	}
}

// consolidate put the values of the passed piece in the pit
func (p Pit) consolidate(pc *Piece) {
	for i, tile := range pc.Tiles {
		if pc.Y-i < 0 {
			return
		}
		p[pc.Y-i][pc.X] = tile
	}
}
