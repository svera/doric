package columns

const (
	Remove = -1
	Empty  = 0
)

type Pit struct {
	width  int
	height int
	cells  [][]int
}

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

func (p *Pit) CheckLines() int {
	remove := map[Coords]struct{}{}
	p.checkHorizontalLines(remove)
	p.checkVerticalLines(remove)
	for coords := range remove {
		// Cells with negative values are cells with tiles to be removed
		p.cells[coords.y][coords.x] = Remove
	}
	return len(remove)
}

func (p *Pit) checkHorizontalLines(remove map[Coords]struct{}) {
	for y := p.height - 1; y >= 0; y-- {
		x := 0
		for x < p.width-2 {
			if p.cells[y][x] == Empty {
				x++
				continue
			}
			if p.cells[y][x] == p.cells[y][x+1] && p.cells[y][x+1] == p.cells[y][x+2] {
				remove[Coords{x, y}] = struct{}{}
				remove[Coords{x + 1, y}] = struct{}{}
				remove[Coords{x + 2, y}] = struct{}{}
			}
			x++
		}
	}
}

func (p *Pit) checkVerticalLines(remove map[Coords]struct{}) {
	for x := 0; x < p.width; x++ {
		y := p.height - 1
		for y > 1 {
			if p.cells[y][x] == Empty {
				break
			}
			if p.cells[y][x] == p.cells[y-1][x] && p.cells[y-1][x] == p.cells[y-2][x] {
				remove[Coords{x, y}] = struct{}{}
				remove[Coords{x, y - 1}] = struct{}{}
				remove[Coords{x, y - 2}] = struct{}{}
			}
			y--
		}
	}
}

func (p *Pit) Cell(x, y int) int {
	return p.cells[y][x]
}

func (p *Pit) Width() int {
	return p.width
}

func (p *Pit) Height() int {
	return p.height
}

// Settle move down all tiles which have empty cells below
func (p *Pit) Settle() {
	for x := 0; x < p.Width(); x++ {
		tiles := []int{}
		for y := p.Height() - 1; y >= 0; y-- {
			// This cell contains a tile to be removed, do not put it in the slice of tiles to move settle
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

func (p *Pit) Consolidate(pc *Piece) {
	for i, tile := range pc.Tiles() {
		if pc.Y()-i < 0 {
			return
		}
		p.cells[pc.Y()-i][pc.X()] = tile
	}
}
