package doric

import (
	"time"
)

// Possible commands coming from the player
const (
	CommandLeft = iota
	CommandRight
	CommandDown
	CommandRotate
	CommandPause
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

// Play starts the game loop in a separate thread, making pieces fall to the bottom of the pit at gradually quicker speeds
// as level increases.
// Game can be controlled sending action codes to the input channel. Game updates are communicated as events in the returned
// channel.
// Game ends when no more new pieces can enter the pit, and this will be signaled with the closing of the
// events channel.
func Play(p Pit, rand Randomizer, cfg Config, input <-chan int) <-chan interface{} {
	events := make(chan interface{})

	go func() {
		ticker := time.NewTicker(cfg.Frequency)
		ticks := 0
		totalRemoved := 0

		defer func() {
			close(events)
		}()

		pit := NewPit(p.height(), p.width())
		copy(pit, p)
		current := NewPiece(rand)
		next := NewPiece(rand)
		current.X = pit.width() / 2
		combo := 1
		slowdown := cfg.InitialSlowdown
		level := 1
		paused := false

		sendEventRenewed(events, pit, current, next)
		for {
			select {
			case act := <-input:
				switch act {
				case CommandLeft:
					current.left(pit)
				case CommandRight:
					current.right(pit)
				case CommandDown:
					current.down(pit)
					ticks = 0
				case CommandRotate:
					current.rotate()
				case CommandPause:
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
				if current.down(pit) {
					sendEventUpdated(events, current, paused)
					continue
				}
				pit.consolidate(current)
				removed := pit.markTilesToRemove()
				for removed > 0 {
					totalRemoved += removed
					if slowdown > 1 {
						slowdown--
					}
					if totalRemoved/cfg.NumberTilesForNextLevel > level-1 {
						level++
					}
					sendEventScored(events, pit, removed, combo, level)
					combo++
					pit.settle()
					removed = pit.markTilesToRemove()
				}
				combo = 1
				current.copy(next, pit.width()/2)
				next.randomize(rand)
				sendEventRenewed(events, pit, current, next)

				if pit.Cell(pit.width()/2, 0) != Empty {
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
	p := NewPit(pit.height(), pit.width())
	for i := range pit {
		copy(p[i], pit[i])
	}
	events <- EventScored{
		Pit:     p,
		Combo:   combo,
		Level:   level,
		Removed: total,
	}
}

func sendEventRenewed(events chan<- interface{}, pit Pit, current, next *Piece) {
	p := NewPit(pit.height(), pit.width())
	for i := range pit {
		copy(p[i], pit[i])
	}
	events <- EventRenewed{
		Pit:     p,
		Current: *current,
		Next:    *next,
	}
}
