package columns

import (
	"time"
)

// Events thrown by the game
const (
	EventUpdated = iota
	EventScored
	EventRenewed
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
	// How many points each destroyed tile awards the player
	PointsPerTile int
	// How many tiles a player has to destroy to advance to the next level
	NumberTilesForNextLevel int
	// As the game loop running frequency is every 200ms, an initialSlowdown of 8 means that pieces fall
	// at a speed of 10*200 = 0.5 cells/sec
	// For an updating frequency of 200ms, the maximum falling speed would be 5 cells/sec (a cell every 200ms)
	InitialSlowdown int
	// Frequency to check for tiles to remove, piece changing, etc.
	Frequency time.Duration
}

type status struct {
	Current Piece
	Next    Piece
	Pit     Pit
	Points  int
	Combo   int
	Level   int
	Paused  bool
}

// Event contains the status of the game to be consumed by a client
type Event struct {
	ID     int
	Status status
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
func NewGame(p Pit, r Randomizer, cfg Config) *Game {
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
	g.current.x = p.Width() / 2
	return g
}

// Play starts the game loop, making pieces fall to the bottom of the pit at gradually quicker speeds
// as level increases. Game ends when no more new pieces can enter the pit.
func (g *Game) Play(input <-chan int, events chan<- Event) {
	ticker := time.NewTicker(g.cfg.Frequency)
	ticks := 0
	totalRemoved := 0

	defer func() {
		close(events)
	}()

	g.sendUpdate(events, EventUpdated)
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
			g.sendUpdate(events, EventUpdated)
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
				g.sendUpdate(events, EventUpdated)
				continue
			}
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
				g.sendUpdate(events, EventScored)
			}
			g.combo = 1
			g.current.copy(g.next, g.pit.Width()/2)
			g.next.randomize(g.rand)
			g.sendUpdate(events, EventRenewed)

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

func (g *Game) sendUpdate(events chan<- Event, eventID int) {
	event := Event{
		ID: eventID,
		Status: status{
			Pit:     NewPit(g.pit.Height(), g.pit.Width()),
			Current: *g.current,
			Next:    *g.next,
			Points:  g.points,
			Combo:   g.combo,
			Level:   g.level,
			Paused:  g.paused,
		},
	}

	copy(event.Status.Pit, g.pit)
	events <- event
}
