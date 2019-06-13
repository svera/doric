package columns

import (
	"time"
)

// Events thrown by the game
const (
	Scored = iota
	Finished
)

const (
	pointsPerTile           = 10
	numberTilesForNextLevel = 10
	// As the game loop running frequency every 200ms, an initialSlowdown of 8 means that pieces fall
	// at a speed of 10*200 = 0.5 cells/sec
	// For a updating frwequency of 200ms, the maximum falling speed would be 5 cells/sec (a cell every 200ms)
	initialSlowdown = 10
	frequency       = 200
)

// Game implements the game flow, keeping track of the game's status for a player
type Game struct {
	current  *Piece
	next     *Piece
	pit      *Pit
	points   int
	combo    int
	slowdown int
	paused   bool
	gameOver bool
	level    int
}

// NewGame returns a new Game instance
func NewGame(pit *Pit) *Game {
	g := &Game{
		pit:     pit,
		current: NewPiece(pit),
		next:    NewPiece(pit),
		level:   1,
	}
	g.Reset()
	return g
}

// Play starts the game loop, making pieces fall to the bottom of the pit at gradually quicker speeds
// as level increases. Game ends when no more new pieces can enter the pit.
func (g *Game) Play(events chan<- int) {
	ticker := time.NewTicker(frequency * time.Millisecond)
	ticks := 0
	totalRemoved := 0
	for range ticker.C {
		if g.paused {
			continue
		}
		if ticks != g.slowdown {
			ticks++
			continue
		}
		ticks = 0
		if !g.current.Down() {
			g.pit.consolidate(g.current)
			removed := g.pit.checkLines()
			for removed > 0 {
				totalRemoved += removed
				g.pit.settle()
				g.points += removed * g.combo * pointsPerTile
				g.combo++
				events <- Scored
				removed = g.pit.checkLines()
				if g.slowdown > 1 && totalRemoved/numberTilesForNextLevel > g.level-1 {
					g.slowdown--
					g.level++
				}
			}
			g.combo = 1
			g.current.Copy(g.next)
			g.next.Randomize()
			if g.pit.Cell(g.pit.width/2, 0) != Empty {
				ticker.Stop()
				g.gameOver = true
				events <- Finished
				return
			}
		}
	}
}

// Score returns player's current score
func (g *Game) Score() int {
	return g.points
}

// Level returns player's current level
func (g *Game) Level() int {
	return g.level
}

// Current returns player's current piece falling
func (g *Game) Current() *Piece {
	return g.current
}

// Next returns player's next piece to be played
func (g *Game) Next() *Piece {
	return g.next
}

// Pit returns player's pit
func (g *Game) Pit() *Pit {
	return g.pit
}

// Pause stops game until executed again
func (g *Game) Pause() {
	g.paused = !g.paused
}

// IsPaused returns true if the game is paused
func (g *Game) IsPaused() bool {
	return g.paused
}

// IsGameOver returns true if the game is over
func (g *Game) IsGameOver() bool {
	return g.gameOver
}

// Reset empties pit and reset all game properties to its initial values
func (g *Game) Reset() {
	g.pit.reset()
	g.combo = 1
	g.slowdown = initialSlowdown
	g.points = 0
	g.current.Reset()
	g.next.Reset()
	g.paused = false
	g.gameOver = false
	g.level = 1
}
