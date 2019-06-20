package columns

import (
	"testing"

	"github.com/svera/doric/pkg/columns/mocks"
)

func TestMovement(t *testing.T) {
	r := &mocks.Randomizer{Values: []int{1}}
	p := NewPiece(NewPit(13, 6), r)
	p.tiles = [3]int{1, 2, 3}
	p.x = 0
	p.Left()
	if p.x != 0 {
		t.Errorf("Piece should not move to the left if is in pit's first column")
	}
	p.x = 2
	p.Left()
	if p.x != 1 {
		t.Errorf("Piece should move to the left if isn't in pit's first column")
	}
	p.x = 2
	p.pit.Cells[0][1] = 1
	p.Left()
	if p.x != 2 {
		t.Errorf("Piece should not move to the left if that pit cell is not empty")
	}
	p.x = 5
	p.Right()
	if p.x != 5 {
		t.Errorf("Piece should not move to the right if is in pit's last column")
	}
	p.x = 2
	p.Right()
	if p.x != 3 {
		t.Errorf("Piece should move to the right if isn't in pit's last column")
	}
	p.x = 2
	p.pit.Cells[0][3] = 1
	p.Right()
	if p.x != 2 {
		t.Errorf("Piece should not move to the right if that pit cell is not empty")
	}
	p.y = 0
	p.Down()
	if p.y != 1 {
		t.Errorf("Piece should move down if isn't in pit's bottom")
	}
	p.y = 0
	p.x = 2
	p.pit.Cells[1][2] = 1
	p.Down()
	if p.y != 0 {
		t.Errorf("Piece should not move down if cell pit below is not empty")
	}
}
