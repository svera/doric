package columns

import (
	"reflect"
	"testing"
)

func TestCheckLines(t *testing.T) {
	p := NewPit(13, 6)
	p.cells = [][]int{
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 1, 0, 0, 1, 0},
		[]int{1, 1, 1, 1, 1, 1},
		[]int{1, 1, 0, 1, 1, 1},
	}
	expected := [][]int{
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, -1, 0, 0, -1, 0},
		[]int{-1, -1, -1, -1, -1, -1},
		[]int{1, -1, 0, -1, -1, -1},
	}
	p.checkLines()
	if !reflect.DeepEqual(p.cells, expected) {
		t.Errorf("Expected %v, got %v", expected, p.cells)
	}
}

func TestCheckDiagonalLines(t *testing.T) {
	p := NewPit(13, 6)
	p.cells = [][]int{
		[]int{1, 0, 0, 0, 0, 1},
		[]int{0, 1, 0, 0, 1, 0},
		[]int{0, 0, 1, 1, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{1, 0, 0, 0, 0, 1},
		[]int{0, 1, 0, 0, 1, 0},
		[]int{0, 0, 1, 1, 0, 0},
	}
	expected := [][]int{
		[]int{-1, 0, 0, 0, 0, -1},
		[]int{0, -1, 0, 0, -1, 0},
		[]int{0, 0, -1, -1, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{-1, 0, 0, 0, 0, -1},
		[]int{0, -1, 0, 0, -1, 0},
		[]int{0, 0, -1, -1, 0, 0},
	}
	p.checkLines()
	if !reflect.DeepEqual(p.cells, expected) {
		t.Errorf("Expected %v, got %v", expected, p.cells)
	}
}

func TestSettle(t *testing.T) {
	p := NewPit(13, 6)
	p.cells = [][]int{
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 1, 2, 0, 0, 0},
		[]int{-1, -1, -1, 0, 0, 1},
		[]int{1, 2, 3, 1, 4, 1},
	}
	expected := [][]int{
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0},
		[]int{0, 1, 2, 0, 0, 1},
		[]int{1, 2, 3, 1, 4, 1},
	}
	p.settle()
	if !reflect.DeepEqual(p.cells, expected) {
		t.Errorf("Expected %v, got %v", expected, p.cells)
	}
}

func TestCell(t *testing.T) {
	p := NewPit(13, 6)
	p.cells[12][0] = 1
	p.cells[0][5] = 1
	if p.Cell(0, 12) != 1 || p.Cell(5, 0) != 1 {
		t.Errorf("Cell() not returning the right value")
	}
}
