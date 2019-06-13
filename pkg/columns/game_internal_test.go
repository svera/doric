package columns

import (
	"testing"
)

const (
	pithWidth = 6
	pitHeight = 13
)

func TestGameOver(t *testing.T) {
	pit := NewPit(pitHeight, pithWidth)
	game := NewGame(pit)
	events := make(chan int)
	pit.cells[0][3] = 1
	game.current.y = 12
	go game.Play(events)
	select {
	case ev := <-events:
		if ev == Finished {
			return
		}
	}
}

func TestScore(t *testing.T) {
	pit := NewPit(pitHeight, pithWidth)
	game := NewGame(pit)
	events := make(chan int)
	pit.cells[12][0] = 1
	pit.cells[12][1] = 1
	pit.cells[12][2] = 1
	game.current.x = 5
	game.current.y = 12
	go game.Play(events)

	select {
	case ev := <-events:
		if ev == Scored {
			if game.Score() != 30 {
				t.Errorf("Score should be 30, got %d", game.Score())
			}

			return
		}
	}

}

func TestLevel(t *testing.T) {
	pit := NewPit(pitHeight, pithWidth)
	game := NewGame(pit)
	if game.Level() != 1 {
		t.Errorf("Level should be 1, got %d", game.Level())
	}
}

func TestCurrent(t *testing.T) {
	pit := NewPit(pitHeight, pithWidth)
	game := NewGame(pit)
	p := &Piece{
		tiles: []int{1, 2, 3},
	}
	game.current = p
	if game.Current() != p {
		t.Errorf("Current piece not returned")
	}
}

func TestNext(t *testing.T) {
	pit := NewPit(pitHeight, pithWidth)
	game := NewGame(pit)
	p := &Piece{
		tiles: []int{1, 2, 3},
	}
	game.next = p
	if game.Next() != p {
		t.Errorf("Next piece not returned")
	}
}

func TestPit(t *testing.T) {
	pit := NewPit(pitHeight, pithWidth)
	game := NewGame(pit)
	if game.Pit() != pit {
		t.Errorf("Pit not returned")
	}
}

func TestPause(t *testing.T) {
	pit := NewPit(pitHeight, pithWidth)
	game := NewGame(pit)
	if game.IsPaused() {
		t.Errorf("Game shouldn't be in paused state")
	}
	game.Pause()
	if !game.IsPaused() {
		t.Errorf("Game should be in paused state")
	}
}

func TestIsGameOver(t *testing.T) {
	pit := NewPit(pitHeight, pithWidth)
	game := NewGame(pit)
	if game.IsGameOver() {
		t.Errorf("Game shouldn't be over")
	}
	game.gameOver = true
	if !game.IsGameOver() {
		t.Errorf("Game should be over")
	}
}
