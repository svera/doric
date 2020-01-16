package columns

import (
	"time"
)

// Events thrown by the game
const (
	StatusUpdated = iota
	StatusScored
	StatusRenewed
	StatusPaused
	StatusFinished
)

// Possible actions coming from the player
const (
	ActionLeft = iota
	ActionRight
	ActionDown
	ActionRotate
	ActionReset
	ActionPause
)

// Config holds different parameters related with the game
type Config struct {
	PointsPerTile           int
	NumberTilesForNextLevel int
	// As the game loop running frequency is every 200ms, an initialSlowdown of 8 means that pieces fall
	// at a speed of 10*200 = 0.5 cells/sec
	// For an updating frequency of 200ms, the maximum falling speed would be 5 cells/sec (a cell every 200ms)
	InitialSlowdown int
	Frequency       time.Duration
}

// Update contains the status of the game to be consumed by a client
type Update struct {
	Current Piece
	Next    Piece
	Pit     Pit
	Points  int
	Combo   int
	Status  int
	Level   int
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
}

// NewGame returns a new Game instance
func NewGame(p Pit, current Piece, next Piece, r Randomizer, cfg Config) *Game {
	g := &Game{
		pit:     p,
		current: &current,
		next:    &next,
		level:   1,
		rand:    r,
		cfg:     cfg,
	}
	g.Reset()
	return g
}

// Play starts the game loop, making pieces fall to the bottom of the pit at gradually quicker speeds
// as level increases. Game ends when no more new pieces can enter the pit.
func (g *Game) Play(input <-chan int, updates chan<- Update) {
	ticker := time.NewTicker(g.cfg.Frequency)
	ticks := 0
	totalRemoved := 0

	defer func() {
		close(updates)
	}()

	for {
		select {
		case act := <-input:
			status := StatusUpdated
			if act == ActionLeft {
				g.current.Left()
			}
			if act == ActionRight {
				g.current.Right()
			}
			if act == ActionDown {
				g.current.Down()
			}
			if act == ActionRotate {
				g.current.Rotate()
			}
			if act == ActionPause {
				g.pause()
				if g.paused {
					status = StatusPaused
				}
			}
			g.sendUpdate(updates, status)
		case <-ticker.C:
			if g.paused {
				g.sendUpdate(updates, StatusPaused)
				continue
			}
			if ticks != g.slowdown {
				ticks++
				g.sendUpdate(updates, StatusUpdated)
				continue
			}
			ticks = 0
			if !g.current.Down() {
				g.pit.consolidate(g.current)
				removed := g.pit.markTilesToRemove()
				for removed > 0 {
					totalRemoved += removed
					g.pit.settle()
					g.points += removed * g.combo * g.cfg.PointsPerTile
					g.combo++
					removed = g.pit.markTilesToRemove()
					if g.slowdown > 1 {
						g.slowdown--
					}
					if totalRemoved/g.cfg.NumberTilesForNextLevel > g.level-1 {
						g.level++
					}
					g.sendUpdate(updates, StatusScored)
				}
				g.combo = 1
				g.current.copy(g.next)
				g.next.randomize(g.rand)
				g.sendUpdate(updates, StatusRenewed)

				if g.pit.Cell(g.pit.Width()/2, 0) != Empty {
					ticker.Stop()
					g.sendUpdate(updates, StatusFinished)
					return
				}
			}
		}
	}
}

// pause switch game between playing and paused
func (g *Game) pause() {
	g.paused = !g.paused
}

// Reset empties pit and reset all game properties to its initial values
func (g *Game) Reset() {
	g.pit.reset()
	g.combo = 1
	g.slowdown = g.cfg.InitialSlowdown
	g.points = 0
	g.current.reset(g.rand)
	g.next.randomize(g.rand)
	g.paused = false
	g.level = 1
}

func (g *Game) sendUpdate(updates chan<- Update, status int) {
	updates <- Update{
		Current: *g.current,
		Next:    *g.next,
		Pit:     g.pit,
		Points:  g.points,
		Combo:   g.combo,
		Status:  status,
		Level:   g.level,
	}
}
