package columns_test

import (
	"reflect"
	"testing"

	"github.com/svera/doric/pkg/columns"
	"github.com/svera/doric/pkg/columns/mocks"
)

func TestRotate(t *testing.T) {
	r := &mocks.Randomizer{Values: []int{0, 1, 2}}
	p := columns.NewPiece(columns.NewPit(13, 6), r)
	expected := [3]int{3, 1, 2}
	p.Rotate()
	if !reflect.DeepEqual(p.Tiles(), expected) {
		t.Errorf("Expected %v, got %v", expected, p.Tiles())
	}
}

func TestCopy(t *testing.T) {
	r := &mocks.Randomizer{Values: []int{0, 1, 2, 3, 4, 5}}
	pit := columns.NewPit(13, 6)
	p1 := columns.NewPiece(pit, r)
	p2 := columns.NewPiece(pit, r)
	p1.Down()
	p1.Right()
	p2.Copy(p1)
	if !reflect.DeepEqual(p1.Tiles(), p2.Tiles()) {
		t.Errorf("p2 tiles are not equal as p1")
	}
	if p2.X() != 3 {
		t.Errorf("p2 X coordinate not reset after copy")
	}
	if p2.Y() != 0 {
		t.Errorf("p2 Y coordinate not reset after copy")
	}
}
