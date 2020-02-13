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
	// Pause / unpause game (intended for internal use, e. g. stop game logic while play animation)
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

type game struct {
	pit          Pit
	current      *Piece
	next         *Piece
	combo        int
	slowdown     int
	level        int
	paused       bool
	wait         bool
	ticks        int
	totalRemoved int
	cfg          Config
	events       chan interface{}
	ticker       *time.Ticker
}

// Play starts the game loop in a separate thread, making pieces fall to the bottom of the pit at gradually quicker speeds
// as level increases.
// Game can be controlled sending command codes to the commands channel. Game updates are communicated as events in the returned
// channel.
// Game ends when no more new pieces can enter the pit, and this will be signaled with the closing of the
// events channel.
func Play(p Pit, rand Randomizer, cfg Config, commands <-chan int) <-chan interface{} {
	game := newGame(p, rand, cfg)

	go func() {
		defer func() {
			close(game.events)
			game.ticker.Stop()
		}()

		game.renewPieces(rand)
		for {
			select {
			case comm := <-commands:
				if comm == CommandQuit {
					return
				}
				game.execute(comm)
			case <-game.ticker.C:
				if game.paused || game.wait {
					continue
				}
				if game.ticks != game.slowdown {
					game.ticks++
					continue
				}
				game.ticks = 0
				if game.current.down(game.pit) {
					game.events <- EventUpdated{
						Current: *game.current,
					}
					continue
				}
				game.pit.lock(game.current)
				game.removeLines()
				game.renewPieces(rand)

				if game.pit[game.pit.width()/2][0] != Empty {
					return
				}
			}
		}
	}()

	return game.events
}

func newGame(p Pit, rand Randomizer, cfg Config) *game {
	game := &game{
		pit:      p.copy(),
		current:  &Piece{Tiles: [3]int{}},
		next:     &Piece{Tiles: [3]int{}},
		combo:    1,
		slowdown: cfg.InitialSlowdown,
		level:    1,
		cfg:      cfg,
		events:   make(chan interface{}),
		ticker:   time.NewTicker(cfg.Frequency),
	}
	game.next.randomize(rand)
	return game
}

func (g *game) execute(comm int) {
	if comm == CommandWaitSwitch {
		g.wait = !g.wait
		return
	}
	if comm == CommandPauseSwitch {
		g.paused = !g.paused
		return
	}
	if !g.paused && !g.wait {
		switch comm {
		case CommandLeft:
			g.current.left(g.pit)
		case CommandRight:
			g.current.right(g.pit)
		case CommandDown:
			g.current.down(g.pit)
			g.ticks = 0
		case CommandRotate:
			g.current.rotate()
		}
	}
	g.events <- EventUpdated{
		Current: *g.current,
	}
}

func (g *game) removeLines() {
	removed := g.pit.markTilesToRemove()
	for removed > 0 {
		g.totalRemoved += removed
		if g.slowdown > 1 {
			g.slowdown--
		}
		if g.totalRemoved/g.cfg.NumberTilesForNextLevel > g.level-1 {
			g.level++
		}
		g.events <- EventScored{
			Pit:     g.pit.copy(),
			Combo:   g.combo,
			Level:   g.level,
			Removed: removed,
		}
		g.combo++
		g.pit.settle()
		removed = g.pit.markTilesToRemove()
	}
}

func (g *game) renewPieces(rand Randomizer) {
	g.current.copy(g.next, g.pit.width()/2)
	g.next.randomize(rand)
	g.combo = 1

	g.events <- EventRenewed{
		Pit:     g.pit.copy(),
		Current: *g.current,
		Next:    *g.next,
	}
}
