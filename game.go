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
	// InitialSpeed is the falling speed at the beginning of the game in cells/second
	InitialSpeed float64
	// SpeedIncrement is how much the speed increases each level in cells/second
	SpeedIncrement float64
	// MaxSpeed is the maximum speed falling pieces can reach
	MaxSpeed float64
}

type game struct {
	well         Well
	current      *Piece
	next         *Piece
	combo        int
	level        int
	paused       bool
	wait         bool
	totalRemoved int
	cfg          Config
	events       chan interface{}
	speed        float64
	maxFrequency time.Duration
	ticker       *time.Ticker
}

// Play starts the game loop in a separate thread, making pieces fall to the bottom of the well at gradually quicker speeds
// as level increases.
// Game can be controlled sending command codes to the commands channel. Game updates are communicated as events in the returned
// channel.
// Game ends when no more new pieces can enter the well, and this will be signaled with the closing of the
// events channel.
func Play(p Well, rand Randomizer, cfg Config, commands <-chan int) <-chan interface{} {
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
				if game.current.down(game.well) {
					game.events <- EventUpdated{
						Current: *game.current,
					}
					continue
				}
				game.well.lock(game.current)
				game.removeLines()
				game.renewPieces(rand)

				if game.isOver() {
					return
				}
			}
		}
	}()

	return game.events
}

func newGame(p Well, rand Randomizer, cfg Config) *game {
	game := &game{
		well:         p.copy(),
		current:      &Piece{Tiles: [3]int{}},
		next:         &Piece{Tiles: [3]int{}},
		combo:        1,
		level:        1,
		cfg:          cfg,
		events:       make(chan interface{}),
		speed:        cfg.InitialSpeed,
		maxFrequency: time.Duration(1000/cfg.MaxSpeed) * time.Millisecond,
		ticker:       time.NewTicker(time.Duration(1000/cfg.InitialSpeed) * time.Millisecond),
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
			g.current.left(g.well)
		case CommandRight:
			g.current.right(g.well)
		case CommandDown:
			g.current.down(g.well)
		case CommandRotate:
			g.current.rotate()
		}
	}
	g.events <- EventUpdated{
		Current: *g.current,
	}
}

func (g *game) removeLines() {
	removed := g.well.markTilesToRemove()
	for removed > 0 {
		g.totalRemoved += removed
		if g.totalRemoved/g.cfg.NumberTilesForNextLevel > g.level-1 {
			g.level++
			g.speedUp()
		}
		g.events <- EventScored{
			Well:    g.well.copy(),
			Combo:   g.combo,
			Level:   g.level,
			Removed: removed,
		}
		g.combo++
		g.well.settle()
		removed = g.well.markTilesToRemove()
	}
}

func (g *game) speedUp() {
	speed := g.speed + g.cfg.SpeedIncrement
	freq := time.Duration(1000/speed) * time.Millisecond
	if freq > g.maxFrequency {
		g.ticker.Stop()
		g.speed = speed
		g.ticker = time.NewTicker(freq)
	}
}

func (g *game) isOver() bool {
	return g.well[g.well.width()/2][0] != Empty
}

func (g *game) renewPieces(rand Randomizer) {
	g.current.copy(g.next, g.well.width()/2)
	g.next.randomize(rand)
	g.combo = 1

	g.events <- EventRenewed{
		Well:    g.well.copy(),
		Current: *g.current,
		Next:    *g.next,
	}
}
