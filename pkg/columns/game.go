package columns

import (
	"time"
)

// Possible actions coming from the player
const (
	ActionLeft = iota
	ActionRight
	ActionDown
	ActionRotate
	ActionPause
)

// Config holds different parameters related with the game
type Config struct {
	// How many tiles a player has to destroy to advance to the next level
	NumberTilesForNextLevel int
	// As the game loop running frequency is every 200ms, an initialSlowdown of 8 means that pieces fall
	// at a speed of 10*200 = 0.5 cells/sec
	// For an updating frequency of 200ms, the maximum falling speed would be 5 cells/sec (a cell every 200ms)
	InitialSlowdown int
	// Frequency to check for tiles to remove, piece changing, etc.
	Frequency time.Duration
}

// Play starts the game loop, making pieces fall to the bottom of the pit at gradually quicker speeds
// as level increases. Game ends when no more new pieces can enter the pit.
func Play(pit Pit, rand Randomizer, cfg Config, input <-chan int) <-chan interface{} {
	events := make(chan interface{})

	go func() {
		ticker := time.NewTicker(cfg.Frequency)
		ticks := 0
		totalRemoved := 0

		defer func() {
			close(events)
		}()

		current := NewPiece(rand)
		next := NewPiece(rand)
		current.X = pit.Width() / 2
		combo := 1
		slowdown := cfg.InitialSlowdown
		level := 1
		paused := false

		sendEventRenewed(events, current, next)
		for {
			select {
			case act := <-input:
				switch act {
				case ActionLeft:
					current.Left(pit)
				case ActionRight:
					current.Right(pit)
				case ActionDown:
					current.Down(pit)
				case ActionRotate:
					current.Rotate()
				case ActionPause:
					paused = !paused
				}
				sendEventUpdated(events, current, paused)
			case <-ticker.C:
				if paused {
					continue
				}
				if ticks != slowdown {
					ticks++
					continue
				}
				ticks = 0
				if current.Down(pit) {
					sendEventUpdated(events, current, paused)
					continue
				}
				pit.consolidate(current)
				removed := pit.markTilesToRemove()
				for removed > 0 {
					totalRemoved += removed
					pit.settle()
					if slowdown > 1 {
						slowdown--
					}
					if totalRemoved/cfg.NumberTilesForNextLevel > level-1 {
						level++
					}
					sendEventScored(events, pit, removed, combo, level)
					combo++
					removed = pit.markTilesToRemove()
				}
				combo = 1
				current.copy(next, pit.Width()/2)
				next.randomize(rand)
				sendEventRenewed(events, current, next)

				if pit.Cell(pit.Width()/2, 0) != Empty {
					ticker.Stop()
					return
				}
			}
		}
	}()

	return events
}

func sendEventUpdated(events chan<- interface{}, current *Piece, paused bool) {
	events <- EventUpdated{
		Current: *current,
		Paused:  paused,
	}
}

func sendEventScored(events chan<- interface{}, pit Pit, total, combo, level int) {
	p := NewPit(pit.Height(), pit.Width())
	copy(p, pit)
	events <- EventScored{
		Pit:     p,
		Combo:   combo,
		Level:   level,
		Removed: total,
	}
}

func sendEventRenewed(events chan<- interface{}, current, next *Piece) {
	events <- EventRenewed{
		Current: *current,
		Next:    *next,
	}
}
