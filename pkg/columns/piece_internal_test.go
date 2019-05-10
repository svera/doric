package columns

import (
	"reflect"
	"testing"
)

func TestRotate(t *testing.T) {
	p := NewPiece(NewPit(13, 6))
	p.tiles = []int{1, 2, 3}
	expected := []int{3, 1, 2}
	p.Rotate()
	if !reflect.DeepEqual(p.tiles, expected) {
		t.Errorf("Expected %v, got %v", expected, p.tiles)
	}
}
