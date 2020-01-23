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
	current  *Piece
	next     *Piece
	pit      Pit
	points   int
	combo    int
	slowdown int
	paused   bool
	level    int
	rand     Randomizer
	cfg      Config
	events   chan interface{}
}

// NewGame returns a new Game instance
func NewGame(p Pit, r Randomizer, cfg Config) (*Game, <-chan interface{}) {
	g := &Game{
		current:  NewPiece(r),
		next:     NewPiece(r),
		pit:      p,
		combo:    1,
		slowdown: cfg.InitialSlowdown,
		level:    1,
		rand:     r,
		cfg:      cfg,
	}
	g.current.X = p.Width() / 2
	g.events = make(chan interface{})
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

	g.sendEventRenewed()
	for {
		select {
		case act := <-input:
			switch act {
			case ActionLeft:
				g.current.Left(g.pit)
			case ActionRight:
				g.current.Right(g.pit)
			case ActionDown:
				g.current.Down(g.pit)
			case ActionRotate:
				g.current.Rotate()
			case ActionPause:
				g.pause()
			}
			g.sendEventUpdated()
		case <-ticker.C:
			if g.paused {
				continue
			}
			if ticks != g.slowdown {
				ticks++
				continue
			}
			ticks = 0
			if g.current.Down(g.pit) {
				g.sendEventUpdated()
				continue
			}
			g.pit.consolidate(g.current)
			removed := g.pit.markTilesToRemove()
			for removed > 0 {
				totalRemoved += removed
				g.pit.settle()
				if g.slowdown > 1 {
					g.slowdown--
				}
				if totalRemoved/g.cfg.NumberTilesForNextLevel > g.level-1 {
					g.level++
				}
				g.sendEventScored(removed)
				g.combo++
				removed = g.pit.markTilesToRemove()
			}
			g.combo = 1
			g.current.copy(g.next, g.pit.Width()/2)
			g.next.randomize(g.rand)
			g.sendEventRenewed()

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

func (g *Game) sendEventUpdated() {
	g.events <- EventUpdated{
		Current: *g.current,
		Paused:  g.paused,
	}
}

func (g *Game) sendEventScored(total int) {
	p := NewPit(g.pit.Height(), g.pit.Width())
	copy(p, g.pit)
	g.events <- EventScored{
		Pit:     p,
		Combo:   g.combo,
		Level:   g.level,
		Removed: total,
	}
}

func (g *Game) sendEventRenewed() {
	g.events <- EventRenewed{
		Current: *g.current,
		Next:    *g.next,
	}
}
