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
)

// Player implements the game flow, keeping track of the game's status for a player
type Player struct {
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

// NewPlayer returns a new Player instance
func NewPlayer(pit *Pit) *Player {
	return &Player{
		pit:     pit,
		current: NewPiece(pit),
		next:    NewPiece(pit),
		level:   1,
	}
}

// Play starts the game loop, making pieces fall to the bottom of the pit at gradually quicker speeds
// as level increases. Game ends when no more new pieces can enter the pit.
func (p *Player) Play(events chan<- int) {
	ticker := time.NewTicker(200 * time.Millisecond)
	ticks := 0
	totalRemoved := 0
	for range ticker.C {
		if p.paused {
			continue
		}
		if ticks != p.slowdown {
			ticks++
			continue
		}
		ticks = 0
		if !p.current.Down() {
			p.pit.consolidate(p.current)
			removed := p.pit.checkLines()
			for removed > 0 {
				totalRemoved += removed
				p.pit.settle()
				p.points += removed * p.combo * pointsPerTile
				p.combo++
				events <- Scored
				removed = p.pit.checkLines()
				if p.slowdown > 1 && totalRemoved/numberTilesForNextLevel > p.level-1 {
					p.slowdown--
					p.level++
				}
			}
			p.combo = 1
			p.current.Copy(p.next)
			p.next.Randomize()
		}
		if p.pit.Cell(p.pit.width/2, 0) != Empty {
			ticker.Stop()
			p.gameOver = true
			events <- Finished
			return
		}
	}
}

// Score returns player's current score
func (p *Player) Score() int {
	return p.points
}

// Level returns player's current level
func (p *Player) Level() int {
	return p.level
}

// Current returns player's current piece falling
func (p *Player) Current() *Piece {
	return p.current
}

// Next returns player's next piece to be played
func (p *Player) Next() *Piece {
	return p.next
}

// Pit returns player's pit
func (p *Player) Pit() *Pit {
	return p.pit
}

// Pause stops game until executed again
func (p *Player) Pause() {
	p.paused = !p.paused
}

// IsPaused returns true if the game is paused
func (p *Player) IsPaused() bool {
	return p.paused
}

// IsGameOver returns true if the game is over
func (p *Player) IsGameOver() bool {
	return p.gameOver
}

// Reset empties pit and reset all game properties to its initial values
func (p *Player) Reset() {
	p.pit.reset()
	p.combo = 1
	p.slowdown = 8
	p.points = 0
	p.current.Reset()
	p.next.Reset()
	p.paused = false
	p.gameOver = false
	p.level = 1
}
