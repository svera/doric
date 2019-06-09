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
	player := NewPlayer(pit)
	events := make(chan int)
	player.Reset()
	pit.cells[0][3] = 1
	go player.Play(events)
	select {
	case ev := <-events:
		if ev == Finished {
			return
		}
	}
}

func TestScore(t *testing.T) {
	pit := NewPit(pitHeight, pithWidth)
	player := NewPlayer(pit)
	events := make(chan int)
	player.Reset()
	pit.cells[12][0] = 1
	pit.cells[12][1] = 1
	pit.cells[12][2] = 1
	player.current.x = 5
	player.current.y = 12
	go player.Play(events)

	select {
	case ev := <-events:
		if ev == Scored {
			if player.Score() != 30 {
				t.Errorf("Score should be 30, got %d", player.Score())
			}

			return
		}
	}

}
