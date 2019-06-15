package columns

import (
	"reflect"
	"testing"
)

func TestRotate(t *testing.T) {
	p := NewPiece(NewPit(13, 6))
	p.tiles = [3]int{1, 2, 3}
	expected := [3]int{3, 1, 2}
	p.Rotate()
	if !reflect.DeepEqual(p.tiles, expected) {
		t.Errorf("Expected %v, got %v", expected, p.tiles)
	}
}

func TestMovement(t *testing.T) {
	p := NewPiece(NewPit(13, 6))
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
	p.pit.cells[0][1] = 1
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
	p.pit.cells[0][3] = 1
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
	p.pit.cells[1][2] = 1
	p.Down()
	if p.y != 0 {
		t.Errorf("Piece should not move down if cell pit below is not empty")
	}
}

func TestCopy(t *testing.T) {
	pit := NewPit(13, 6)
	p1 := NewPiece(pit)
	p2 := NewPiece(pit)
	p1.tiles = [3]int{1, 2, 3}
	p1.x = 1
	p1.y = 1
	p2.Copy(p1)
	if !reflect.DeepEqual(p1.tiles, p2.tiles) {
		t.Errorf("p2 tiles are not equal as p1")
	}
	if p2.x != 3 {
		t.Errorf("p2 X coordinate not reset after copy")
	}
	if p2.y != 0 {
		t.Errorf("p2 Y coordinate not reset after copy")
	}
}
