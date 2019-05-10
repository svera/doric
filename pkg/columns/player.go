package columns

import (
	"time"
)

// Events thrown by the game
const (
	Scored = iota
	Finished
)

const pointsPerTile = 10

// Player implements the game flow, keeping track of the game's status for a player
type Player struct {
	current  *Piece
	next     *Piece
	pit      *Pit
	points   int
	combo    int
	slowdown int
	blocks   int
}

// NewPlayer returns a new Player instance
func NewPlayer(pit *Pit) *Player {
	return &Player{
		pit:      pit,
		combo:    1,
		slowdown: 8,
		current:  NewPiece(pit),
		next:     NewPiece(pit),
	}
}

// Score returns player's current score
func (p *Player) Score() int {
	return p.points
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

func (p *Player) Play(events chan<- *Event) {
	ticker := time.NewTicker(250 * time.Millisecond)
	go func(events chan<- *Event) {
		ticks := 0
		for range ticker.C {
			if ticks != p.slowdown {
				ticks++
				continue
			}
			ticks = 0
			if !p.current.Down() {
				p.pit.Consolidate(p.current)
				removed := p.pit.CheckLines()
				for removed > 0 {
					p.pit.Settle()
					p.points += removed * p.combo * pointsPerTile
					p.combo++
					events <- NewEvent(Scored, struct{}{})
					removed = p.pit.CheckLines()
					if p.slowdown > 1 {
						p.slowdown--
					}
				}
				p.combo = 1
				p.current.Copy(p.next)
				p.next.Renew()
				if p.pit.Cell(3, 0) != Empty {
					ticker.Stop()
					events <- NewEvent(Finished, struct{}{})
					return
				}
			}
		}
	}(events)
}
