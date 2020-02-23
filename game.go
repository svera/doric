package doric

import (
	"fmt"
	"time"
)

// Possible commands coming from the player
const (
	// Move the current column left
	CommandLeft = iota
	// Move the current column right
	CommandRight
	// Move the current column down
	CommandDown
	// Rotate tiles in current column
	CommandRotate
	// Pause / unpause game (for player use)
	CommandPauseSwitch
	// Pause / unpause game (intended for internal use, e. g. stop game logic while play animation)
	CommandWaitSwitch
	// Quit game
	CommandQuit
)

// Possible returned errors
const (
	errorNegativeNumberTilesForNextLevel = "NumberTilesForNextLevel must be equal or greater than 0"
	errorLessEqualZeroInitialSpeed       = "InitialSpeed must be greater than 0"
	errorNegativeSpeedIncrement          = "SpeedIncrement must be equal or greater than 0"
	errorLessEqualZeroMaxSpeed           = "MaxSpeed must be greater than 0"
)

const nanosecond = 1000000000

// Config holds different parameters related with the game
type Config struct {
	// How many tiles a player has to destroy to advance to the next level
	// Must be equal or greater than zero.
	NumberTilesForNextLevel int
	// InitialSpeed is the falling speed at the beginning of the game in cells/second.
	// Must be greater than zero.
	InitialSpeed float64
	// SpeedIncrement is how much the speed increases each level in cells/second
	// Must be equal or greater than zero.
	SpeedIncrement float64
	// MaxSpeed is the maximum speed falling columns can reach
	// Must be greater than zero.
	MaxSpeed float64
}

type game struct {
	well         Well
	column       *Column
	nextTileset  [3]int
	level        int
	paused       bool
	wait         bool
	totalRemoved int
	cfg          Config
	events       chan interface{}
	speed        float64
	ticker       *time.Ticker
	build        TilesetBuilder
}

// Play starts the game loop in a separate thread, making columns fall to the bottom of the well at gradually quicker speeds
// as level increases.
// Game can be controlled sending command codes to the commands channel. Game updates are communicated as events in the returned
// channel.
// Game ends when no more new columns can enter the well, and this will be signaled with the closing of the
// events channel.
func Play(p Well, builder TilesetBuilder, cfg Config, commands <-chan int) (<-chan interface{}, error) {
	game, err := newGame(p, builder, cfg)
	if err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			close(game.events)
			game.ticker.Stop()
		}()

		game.renewColumn()
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
				if game.column.down(game.well) {
					game.events <- EventUpdated{
						Column: *game.column,
					}
					continue
				}
				game.removeLines()
				game.renewColumn()

				if game.isOver() {
					return
				}
			}
		}
	}()

	return game.events, nil
}

func newGame(p Well, build TilesetBuilder, cfg Config) (*game, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return &game{
		well:        p.copy(),
		column:      &Column{Tileset: [3]int{}},
		nextTileset: build(maxTile),
		level:       1,
		cfg:         cfg,
		events:      make(chan interface{}),
		speed:       cfg.InitialSpeed,
		ticker:      time.NewTicker(time.Duration(nanosecond / cfg.InitialSpeed)),
		build:       build,
	}, nil
}

func validateConfig(cfg Config) error {
	if cfg.NumberTilesForNextLevel < 0 {
		return fmt.Errorf(errorNegativeNumberTilesForNextLevel)
	}
	if cfg.InitialSpeed <= 0 {
		return fmt.Errorf(errorLessEqualZeroInitialSpeed)
	}
	if cfg.SpeedIncrement < 0 {
		return fmt.Errorf(errorNegativeSpeedIncrement)
	}
	if cfg.MaxSpeed <= 0 {
		return fmt.Errorf(errorLessEqualZeroMaxSpeed)
	}
	return nil
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
			g.column.left(g.well)
		case CommandRight:
			g.column.right(g.well)
		case CommandDown:
			g.column.down(g.well)
		case CommandRotate:
			g.column.rotate()
		}
	}
	g.events <- EventUpdated{
		Column: *g.column,
	}
}

func (g *game) removeLines() {
	g.well.lock(g.column)
	removed := g.well.markTilesToRemove()
	combo := 1
	for removed > 0 {
		g.totalRemoved += removed
		if g.totalRemoved/g.cfg.NumberTilesForNextLevel > g.level-1 {
			g.level++
			g.speedUp()
		}
		g.events <- EventScored{
			Well:    g.well.copy(),
			Combo:   combo,
			Level:   g.level,
			Removed: removed,
		}
		combo++
		g.well.settle()
		removed = g.well.markTilesToRemove()
	}
}

func (g *game) speedUp() {
	speed := g.speed + g.cfg.SpeedIncrement
	if speed < g.cfg.MaxSpeed {
		g.ticker.Stop()
		g.speed = speed
		g.ticker = time.NewTicker(time.Duration(nanosecond / speed))
	}
}

func (g *game) isOver() bool {
	return g.well[g.well.width()/2][0] != Empty
}

func (g *game) renewColumn() {
	g.column.reset(g.nextTileset, g.well.width()/2)
	g.nextTileset = g.build(maxTile)

	g.events <- EventRenewed{
		Well:        g.well.copy(),
		Column:      *g.column,
		NextTileset: g.nextTileset,
	}
}
