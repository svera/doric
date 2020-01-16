package columns_test

import (
	"testing"
	"time"

	"github.com/svera/doric/pkg/columns"
	"github.com/svera/doric/pkg/columns/mocks"
)

const (
	pithWidth = 6
	pitHeight = 13
)

var codeToStatusName = [6]string{
	"StatusStarted",
	"StatusUpdated",
	"StatusScored",
	"StatusRenewed",
	"StatusPaused",
	"StatusFinished",
}

func TestGameOver(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := columns.NewPit(1, pithWidth)
	current := columns.NewPiece(pit)
	next := columns.NewPiece(pit)
	r := &mocks.Randomizer{Values: []int{0}}
	game := columns.NewGame(pit, *current, *next, r)
	updates := make(chan columns.Update)
	input := make(chan int)
	pit[0][3] = 1
	go game.Play(input, updates)
	select {
	case upd := <-updates:
		if upd.Status == columns.StatusFinished {
			break
		}
	case <-timeout:
		t.Errorf("Game should be over")
	}
}

func TestLevel(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := columns.NewPit(2, pithWidth)
	current := columns.NewPiece(pit)
	next := columns.NewPiece(pit)
	r := &mocks.Randomizer{Values: []int{0}}
	game := columns.NewGame(pit, *current, *next, r)
	updates := make(chan columns.Update)
	input := make(chan int)
	pit[1][0] = 1
	pit[1][1] = 1
	pit[1][2] = 1
	pit[1][3] = 1
	pit[1][4] = 1
	pit[1][5] = 1
	pit[0][0] = 1
	pit[0][1] = 1
	pit[0][2] = 1

	go game.Play(input, updates)

	select {
	case upd := <-updates:
		if upd.Status == columns.StatusScored {
			if upd.Level == 2 {
				return
			}
		}
	case <-timeout:
		t.Errorf("Level should be 2")
	}
}

func TestScore(t *testing.T) {
	timeout := time.After(3 * time.Second)
	pit := columns.NewPit(3, pithWidth)
	current := columns.NewPiece(pit)
	next := columns.NewPiece(pit)
	r := &mocks.Randomizer{Values: []int{0, 1, 2, 3, 4, 5}}
	game := columns.NewGame(pit, *current, *next, r)
	updates := make(chan columns.Update)
	input := make(chan int)
	pit[1][3] = 1
	pit[2][3] = 1
	go game.Play(input, updates)

	select {
	case upd := <-updates:
		if upd.Status == columns.StatusScored {
			if upd.Points != 30 {
				t.Errorf("Score should be 30, got %d", upd.Points)
			}
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	select {
	case upd := <-updates:
		if upd.Status == columns.StatusRenewed {
			expectedTiles := [3]int{4, 5, 6}
			if upd.Current.Tiles() != expectedTiles {
				t.Errorf(
					"Expected that the next piece was copied to the current one with values %v, got %v",
					expectedTiles,
					upd.Current.Tiles(),
				)
			}
			return
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}
}

func TestPause(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := columns.NewPit(pitHeight, pithWidth)
	current := columns.NewPiece(pit)
	next := columns.NewPiece(pit)
	r := &mocks.Randomizer{Values: []int{1}}
	game := columns.NewGame(pit, *current, *next, r)
	updates := make(chan columns.Update)
	input := make(chan int)
	go game.Play(input, updates)

	go func() {
		input <- columns.ActionPause
	}()

	select {
	case upd := <-updates:
		if upd.Status == columns.StatusPaused {
			break
		}
		t.Errorf("Game must be paused, got '%s'", codeToStatusName[upd.Status])
	case <-timeout:
		t.Errorf("Test timed out")
	}

	go func() {
		input <- columns.ActionPause
	}()

	select {
	case upd := <-updates:
		if upd.Status != columns.StatusPaused {
			break
		}
		t.Errorf("Game must not be paused, got '%s'", codeToStatusName[upd.Status])
	case <-timeout:
		t.Errorf("Test timed out")
	}

}

func TestInput(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := columns.NewPit(pitHeight, pithWidth)
	current := columns.NewPiece(pit)
	next := columns.NewPiece(pit)
	r := &mocks.Randomizer{Values: []int{0, 1, 2}}
	game := columns.NewGame(pit, *current, *next, r)
	updates := make(chan columns.Update)
	input := make(chan int)
	go game.Play(input, updates)

	go func() {
		input <- columns.ActionLeft
	}()

	select {
	case upd := <-updates:
		if upd.Current.X() == 2 {
			break
		}
		t.Errorf("Current piece must be at column %d but is at %d", 2, upd.Current.X())
	case <-timeout:
		t.Errorf("Test timed out")
	}

	go func() {
		input <- columns.ActionRight
	}()

	select {
	case upd := <-updates:
		if upd.Current.X() == 3 {
			break
		}
		t.Errorf("Current piece must be at column %d but is at %d", 3, upd.Current.X())
	case <-timeout:
		t.Errorf("Test timed out")
	}

	go func() {
		input <- columns.ActionDown
	}()

	select {
	case upd := <-updates:
		if upd.Current.Y() == 1 {
			break
		}
		t.Errorf("Current piece must be at row %d but is at %d", 1, upd.Current.Y())
	case <-timeout:
		t.Errorf("Test timed out")
	}

	go func() {
		input <- columns.ActionRotate
	}()

	select {
	case upd := <-updates:
		if upd.Current.Tiles() == [3]int{3, 1, 2} {
			break
		}
		t.Errorf("Current piece must be as %v but is as %v", [3]int{3, 1, 2}, upd.Current.Tiles())
	case <-timeout:
		t.Errorf("Test timed out")
	}

}