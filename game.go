package doric

import (
	"time"
)

// Possible commands coming from the player
const (
	// Move the current piece left
	CommandLeft = iota
	// Move the current piece right
	CommandRight
	// Move the current piece down
	CommandDown
	// Rotate tiles in current piece
	CommandRotate
	// Pause / unpause game (for player use)
	CommandPauseSwitch
	// Pause / unpase game (intended for internal use, e. g. stop game logic while play animation)
	CommandWaitSwitch
	// Quit game
	CommandQuit
)

// Config holds different parameters related with the game
type Config struct {
	// How many tiles a player has to destroy to advance to the next level
	NumberTilesForNextLevel int
	// If the game loop running frequency is every 200ms, an initialSlowdown of 8 means that pieces fall
	// at a speed of 10*200 = 0.5 cells/sec
	// For an updating frequency of 200ms, the maximum falling speed would be 5 cells/sec (a cell every 200ms)
	InitialSlowdown int
	// Frequency to check for tiles to remove, piece changing, etc.
	Frequency time.Duration
}

// Play starts the game loop in a separate thread, making pieces fall to the bottom of the pit at gradually quicker speeds
// as level increases.
// Game can be controlled sending command codes to the commands channel. Game updates are communicated as events in the returned
// channel.
// Game ends when no more new pieces can enter the pit, and this will be signaled with the closing of the
// events channel.
func Play(p Pit, rand Randomizer, cfg Config, commands <-chan int) <-chan interface{} {
	events := make(chan interface{})
	pit := NewPit(p.height(), p.width())
	copy(pit, p)
	current := NewPiece(rand)
	next := NewPiece(rand)
	current.X = pit.width() / 2
	combo := 1
	slowdown := cfg.InitialSlowdown
	level := 1
	paused := false
	wait := false
	ticker := time.NewTicker(cfg.Frequency)
	ticks := 0
	totalRemoved := 0

	go func() {
		defer func() {
			close(events)
		}()

		sendEventRenewed(events, pit, current, next)
		for {
			select {
			case comm := <-commands:
				if comm == CommandWaitSwitch {
					wait = !wait
					continue
				}
				if comm == CommandPauseSwitch {
					paused = !paused
					continue
				}
				if comm == CommandQuit {
					return
				}
				if !paused && !wait {
					switch comm {
					case CommandLeft:
						current.left(pit)
					case CommandRight:
						current.right(pit)
					case CommandDown:
						current.down(pit)
						ticks = 0
					case CommandRotate:
						current.rotate()
					}
				}
				sendEventUpdated(events, current)
			case <-ticker.C:
				if paused || wait {
					continue
				}
				if ticks != slowdown {
					ticks++
					continue
				}
				ticks = 0
				if current.down(pit) {
					sendEventUpdated(events, current)
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

				if pit[pit.width()/2][0] != Empty {
					ticker.Stop()
					return
				}
			}
		}
	}()

	return events
}

func sendEventUpdated(events chan<- interface{}, current *Piece) {
	events <- EventUpdated{
		Current: *current,
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
