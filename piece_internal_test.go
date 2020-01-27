package doric

import (
	"reflect"
	"testing"
)

const (
	pithWidth = 6
	pitHeight = 13
)

func TestMovement(t *testing.T) {
	r := &MockRandomizer{Values: []int{0, 1, 2}}
	pit := NewPit(pitHeight, pithWidth)
	p := NewPiece(r)
	p.X = 0
	p.Left(pit)
	if p.X != 0 {
		t.Errorf("Piece should not move to the left if is in pit's first column")
	}
	p.X = 2
	p.Left(pit)
	if p.X != 1 {
		t.Errorf("Piece should move to the left if isn't in pit's first column")
	}
	p.X = 2
	pit[0][1] = 1
	p.Left(pit)
	if p.X != 2 {
		t.Errorf("Piece should not move to the left if that pit cell is not empty")
	}
	p.X = 5
	p.Right(pit)
	if p.X != 5 {
		t.Errorf("Piece should not move to the right if is in pit's last column")
	}
	p.X = 2
	p.Right(pit)
	if p.X != 3 {
		t.Errorf("Piece should move to the right if isn't in pit's last column")
	}
	p.X = 2
	pit[0][3] = 1
	p.Right(pit)
	if p.X != 2 {
		t.Errorf("Piece should not move to the right if that pit cell is not empty")
	}
	p.Y = 0
	p.Down(pit)
	if p.Y != 1 {
		t.Errorf("Piece should move down if isn't in pit's bottom")
	}
	p.Y = 0
	p.X = 2
	pit[1][2] = 1
	p.Down(pit)
	if p.Y != 0 {
		t.Errorf("Piece should not move down if cell pit below is not empty")
	}
}

func TestRotate(t *testing.T) {
	p := &Piece{Tiles: [3]int{1, 2, 3}}
	expected := [3]int{3, 1, 2}
	p.Rotate()
	if !reflect.DeepEqual(p.Tiles, expected) {
		t.Errorf("Expected %v, got %v", expected, p.Tiles)
	}
}
