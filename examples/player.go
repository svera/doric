package main

import (
	"sync"

	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric"
)

// Player handles game's commands in the game, like moving the piece currently falling in the pit
type Player struct {
	*tl.Entity
	Current  *doric.Piece
	Command  chan<- int
	offsetX  int
	offsetY  int
	message  tl.Drawable
	Paused   bool
	Finished bool
	mux      sync.Locker
}

// NewPlayer returns a new Player instance
func NewPlayer(c *doric.Piece, command chan<- int, message tl.Drawable, offsetX, offsetY int, mux sync.Locker) *Player {
	return &Player{
		Current: c,
		Command: command,
		offsetX: offsetX,
		offsetY: offsetY,
		message: message,
		mux:     mux,
	}
}

// Draw draws the piece on screen, as required by Termloop's Drawable interface
// or the paused message if the game is paused
func (p *Player) Draw(screen *tl.Screen) {
	defer p.mux.Unlock()
	p.mux.Lock()
	if p.Paused {
		p.message.(*tl.Text).SetPosition(offsetX+4, offsetY+5)
		p.message.(*tl.Text).SetText("PAUSED")
		return
	}
	if p.Finished {
		p.message.(*tl.Text).SetPosition(offsetX+2, offsetY+5)
		p.message.(*tl.Text).SetText("GAME  OVER")
		return
	}
	p.message.(*tl.Text).SetText("")
	for i := range p.Current.Tiles {
		if i > p.Current.Y {
			continue
		}
		screen.RenderCell(p.Current.X*2+p.offsetX+1, p.Current.Y+p.offsetY-i, &tl.Cell{
			Bg: colors[p.Current.Tiles[i]],
			Fg: tl.ColorBlack,
			Ch: '[',
		})
		screen.RenderCell(p.Current.X*2+p.offsetX+2, p.Current.Y+p.offsetY-i, &tl.Cell{
			Bg: colors[p.Current.Tiles[i]],
			Fg: tl.ColorBlack,
			Ch: ']',
		})
	}
}

// Tick handles events and moves the piece accosdingly if requested, as requested by Termloop's Drawable interface
// as well as the control of the game itself, pausing it
func (p *Player) Tick(event tl.Event) {
	if event.Type == tl.EventKey && !p.Finished { // Is it a keyboard event?
		switch event.Key { // If so, switch on the pressed key.
		case tl.KeyArrowRight:
			p.Command <- doric.CommandRight
		case tl.KeyArrowLeft:
			p.Command <- doric.CommandLeft
		case tl.KeyArrowDown:
			p.Command <- doric.CommandDown
		case tl.KeyTab:
			p.Command <- doric.CommandRotate
		}

		switch event.Ch {
		case 'p', 'P':
			p.Command <- doric.CommandPause
		}
	}
}
