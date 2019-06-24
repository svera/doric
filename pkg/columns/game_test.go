package columns_test

import (
	"testing"

	"github.com/svera/doric/pkg/columns"
	"github.com/svera/doric/pkg/columns/mocks"
)

const (
	pithWidth = 6
	pitHeight = 13
)

func TestGameOver(t *testing.T) {
	pit := columns.NewPit(1, pithWidth)
	r := &mocks.Randomizer{Values: []int{0}}
	game := columns.NewGame(pit, r)
	events := make(chan int)
	pit.Cells[0][3] = 1
	go game.Play(events)
	select {
	case ev := <-events:
		if ev == columns.Finished {
			break
		}
	}
	if !game.IsGameOver() {
		t.Errorf("Game should be over")
	}
}

func TestPit(t *testing.T) {
	pit := columns.NewPit(pitHeight, pithWidth)
	r := &mocks.Randomizer{Values: []int{0}}
	game := columns.NewGame(pit, r)
	if game.Pit() != pit {
		t.Errorf("Pit not returned")
	}
}

func TestLevel(t *testing.T) {
	pit := columns.NewPit(2, pithWidth)
	r := &mocks.Randomizer{Values: []int{0}}
	game := columns.NewGame(pit, r)
	events := make(chan int)
	pit.Cells[1][0] = 1
	pit.Cells[1][1] = 1
	pit.Cells[1][2] = 1
	pit.Cells[1][3] = 1
	pit.Cells[1][4] = 1
	pit.Cells[1][5] = 1
	pit.Cells[0][0] = 1
	pit.Cells[0][1] = 1
	pit.Cells[0][2] = 1

	go game.Play(events)

	select {
	case ev := <-events:
		if ev == columns.Scored {
			if game.Level() == 2 {
				return
			}
		}
	}

	t.Errorf("Level should be 2, got %d", game.Level())
}

func TestScore(t *testing.T) {
	pit := columns.NewPit(3, pithWidth)
	r := &mocks.Randomizer{Values: []int{0, 1, 2, 3, 4, 5}}
	game := columns.NewGame(pit, r)
	events := make(chan int)
	pit.Cells[1][3] = 1
	pit.Cells[2][3] = 1
	go game.Play(events)

	select {
	case ev := <-events:
		if ev == columns.Scored {
			if game.Score() != 30 {
				t.Errorf("Score should be 30, got %d", game.Score())
			}
		}
	}
	select {
	case ev := <-events:
		if ev == columns.Renewed {
			expectedTiles := [3]int{4, 5, 6}
			if game.Current().Tiles() != expectedTiles {
				t.Errorf(
					"Expected that the next piece was copied to the current one with values %v, got %v",
					expectedTiles,
					game.Current().Tiles(),
				)
			}
			return
		}
	}

	t.Errorf("Score event should have been sent and current piece should have been renewed")
}

func TestCurrent(t *testing.T) {
	pit := columns.NewPit(pitHeight, pithWidth)
	r := &mocks.Randomizer{Values: []int{0, 1, 2}}
	game := columns.NewGame(pit, r)
	if game.Current().Tiles() != [3]int{1, 2, 3} {
		t.Errorf("Current piece not returned")
	}
}

func TestNext(t *testing.T) {
	pit := columns.NewPit(pitHeight, pithWidth)
	r := &mocks.Randomizer{Values: []int{5, 5, 5, 0, 1, 2}}
	game := columns.NewGame(pit, r)
	if game.Next().Tiles() != [3]int{1, 2, 3} {
		t.Errorf("Next piece not returned")
	}
}

func TestPause(t *testing.T) {
	pit := columns.NewPit(pitHeight, pithWidth)
	r := &mocks.Randomizer{Values: []int{1}}
	game := columns.NewGame(pit, r)
	if game.IsPaused() {
		t.Errorf("Game shouldn't be in paused state")
	}
	game.Pause()
	if !game.IsPaused() {
		t.Errorf("Game should be in paused state")
	}
}
