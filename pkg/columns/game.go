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

// Game implements the game flow, keeping track of game's status for a player
type Game struct {
	pit    Pit
	paused bool
	rand   Randomizer
	cfg    Config
	events chan interface{}
}

// NewGame returns a new Game instance
func NewGame(p Pit, r Randomizer, cfg Config) (*Game, <-chan interface{}) {
	g := &Game{
		pit:    p,
		rand:   r,
		cfg:    cfg,
		events: make(chan interface{}),
	}
	return g, g.events
}

// Play starts the game loop, making pieces fall to the bottom of the pit at gradually quicker speeds
// as level increases. Game ends when no more new pieces can enter the pit.
func (g *Game) Play(input <-chan int) {
	ticker := time.NewTicker(g.cfg.Frequency)
	ticks := 0
	totalRemoved := 0

	defer func() {
		close(g.events)
	}()

	current := NewPiece(g.rand)
	next := NewPiece(g.rand)
	current.X = g.pit.Width() / 2
	combo := 1
	slowdown := g.cfg.InitialSlowdown
	level := 1

	g.sendEventRenewed(current, next)
	for {
		select {
		case act := <-input:
			switch act {
			case ActionLeft:
				current.Left(g.pit)
			case ActionRight:
				current.Right(g.pit)
			case ActionDown:
				current.Down(g.pit)
			case ActionRotate:
				current.Rotate()
			case ActionPause:
				g.pause()
			}
			g.sendEventUpdated(current)
		case <-ticker.C:
			if g.paused {
				continue
			}
			if ticks != slowdown {
				ticks++
				continue
			}
			ticks = 0
			if current.Down(g.pit) {
				g.sendEventUpdated(current)
				continue
			}
			g.pit.consolidate(current)
			removed := g.pit.markTilesToRemove()
			for removed > 0 {
				totalRemoved += removed
				g.pit.settle()
				if slowdown > 1 {
					slowdown--
				}
				if totalRemoved/g.cfg.NumberTilesForNextLevel > level-1 {
					level++
				}
				g.sendEventScored(removed, combo, level)
				combo++
				removed = g.pit.markTilesToRemove()
			}
			combo = 1
			current.copy(next, g.pit.Width()/2)
			next.randomize(g.rand)
			g.sendEventRenewed(current, next)

			if g.pit.Cell(g.pit.Width()/2, 0) != Empty {
				ticker.Stop()
				return
			}
		}
	}
}

// pause switch game between playing and paused
func (g *Game) pause() {
	g.paused = !g.paused
}

func (g *Game) sendEventUpdated(current *Piece) {
	g.events <- EventUpdated{
		Current: *current,
		Paused:  g.paused,
	}
}

func (g *Game) sendEventScored(total, combo, level int) {
	p := NewPit(g.pit.Height(), g.pit.Width())
	copy(p, g.pit)
	g.events <- EventScored{
		Pit:     p,
		Combo:   combo,
		Level:   level,
		Removed: total,
	}
}

func (g *Game) sendEventRenewed(current, next *Piece) {
	g.events <- EventRenewed{
		Current: *current,
		Next:    *next,
	}
}
