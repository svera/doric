package columns_test

import (
	"testing"

	"github.com/svera/doric/pkg/columns"
)

func TestCell(t *testing.T) {
	p := columns.NewPit(13, 6)
	p[12][0] = 1
	p[0][5] = 1
	if p.Cell(0, 12) != 1 || p.Cell(5, 0) != 1 {
		t.Errorf("Cell() not returning the right value")
	}
}
