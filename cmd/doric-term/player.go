package main

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric/pkg/columns"
)

// Player handles game's interactions in the game, like moving the piece currently falling in the pit
type Player struct {
	*tl.Entity
	game      *columns.Game
	offsetX   int
	offsetY   int
	message   tl.Drawable
	startGame func()
}

// NewPlayer returns a new Player instance
func NewPlayer(p *columns.Game, startGame func(), message tl.Drawable, offsetX int, offsetY int) *Player {
	return &Player{
		game:      p,
		offsetX:   offsetX,
		offsetY:   offsetY,
		message:   message,
		startGame: startGame,
	}
}

// Draw draws the piece on screen, as required by Termloop's Drawable interface
// or the paused message if the game is paused
func (p *Player) Draw(screen *tl.Screen) {
	if p.game.IsPaused() {
		p.message.(*tl.Text).SetPosition(offsetX+4, offsetY+5)
		p.message.(*tl.Text).SetText("PAUSED")
		return
	}
	if p.game.IsGameOver() {
		p.message.(*tl.Text).SetPosition(offsetX+2, offsetY+5)
		p.message.(*tl.Text).SetText("GAME  OVER")
		return
	}
	p.message.(*tl.Text).SetText("")
	for i := range p.game.Current().Tiles() {
		if i > p.game.Current().Y() {
			continue
		}
		screen.RenderCell(p.game.Current().X()*2+p.offsetX+1, p.game.Current().Y()+p.offsetY-i, &tl.Cell{
			Bg: colors[p.game.Current().Tiles()[i]],
			Fg: tl.ColorBlack,
			Ch: '[',
		})
		screen.RenderCell(p.game.Current().X()*2+p.offsetX+2, p.game.Current().Y()+p.offsetY-i, &tl.Cell{
			Bg: colors[p.game.Current().Tiles()[i]],
			Fg: tl.ColorBlack,
			Ch: ']',
		})
	}
}

// Tick handles events and moves the piece accosdingly if requested, as requested by Termloop's Drawable interface
// as well as the control of the game itself, pausing it
func (p *Player) Tick(event tl.Event) {
	if event.Type == tl.EventKey { // Is it a keyboard event?
		switch event.Key { // If so, switch on the pressed key.
		case tl.KeyArrowRight:
			p.game.Current().Right()
		case tl.KeyArrowLeft:
			p.game.Current().Left()
		case tl.KeyArrowDown:
			p.game.Current().Down()
		case tl.KeyTab:
			p.game.Current().Rotate()
		case tl.KeySpace:
			if p.game.IsGameOver() {
				p.game.Reset()
				p.startGame()
			}
		}

		switch event.Ch {
		case 'p', 'P':
			p.game.Pause()
		}
	}
}
