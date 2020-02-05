package main

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric"
)

// Player handles game's interactions in the game, like moving the piece currently falling in the pit
type Player struct {
	*tl.Entity
	Current  *doric.Piece
	Action   chan<- int
	offsetX  int
	offsetY  int
	message  tl.Drawable
	Paused   bool
	Finished bool
}

// NewPlayer returns a new Player instance
func NewPlayer(c *doric.Piece, action chan<- int, message tl.Drawable, offsetX int, offsetY int) *Player {
	return &Player{
		Current: c,
		Action:  action,
		offsetX: offsetX,
		offsetY: offsetY,
		message: message,
	}
}

// Draw draws the piece on screen, as required by Termloop's Drawable interface
// or the paused message if the game is paused
func (p *Player) Draw(screen *tl.Screen) {
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
			p.Action <- doric.CommandRight
		case tl.KeyArrowLeft:
			p.Action <- doric.CommandLeft
		case tl.KeyArrowDown:
			p.Action <- doric.CommandDown
		case tl.KeyTab:
			p.Action <- doric.CommandRotate
		}

		switch event.Ch {
		case 'p', 'P':
			p.Action <- doric.CommandPause
		}
	}
}
