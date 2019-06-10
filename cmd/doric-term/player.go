package main

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric/pkg/columns"
)

// Player handles player's interactions in the game, like moving the piece currently falling in the pit
type Player struct {
	*tl.Entity
	player    *columns.Player
	offsetX   int
	offsetY   int
	message   tl.Drawable
	startGame func()
}

// NewPlayer returns a new Player instance
func NewPlayer(p *columns.Player, startGame func(), message tl.Drawable, offsetX int, offsetY int) *Player {
	return &Player{
		player:    p,
		offsetX:   offsetX,
		offsetY:   offsetY,
		message:   message,
		startGame: startGame,
	}
}

// Draw draws the piece on screen, as required by Termloop's Drawable interface
// or the paused message if the game is paused
func (p *Player) Draw(screen *tl.Screen) {
	if p.player.IsPaused() {
		p.message.(*tl.Text).SetPosition(offsetX+1, offsetY+5)
		p.message.(*tl.Text).SetText("PAUSED")
		return
	}
	if p.player.IsGameOver() {
		p.message.(*tl.Text).SetPosition(offsetX-1, offsetY+5)
		p.message.(*tl.Text).SetText("GAME OVER")
		return
	}
	p.message.(*tl.Text).SetText("")
	for i := range p.player.Current().Tiles() {
		if i > p.player.Current().Y() {
			continue
		}
		screen.RenderCell(p.player.Current().X()+p.offsetX+1, p.player.Current().Y()+p.offsetY-i, &tl.Cell{
			Bg: colors[p.player.Current().Tiles()[i]],
			Fg: tl.ColorBlack,
			Ch: chars[p.player.Current().Tiles()[i]],
		})
	}
}

// Tick handles events and moves the piece accosdingly if requested, as requested by Termloop's Drawable interface
// as well as the control of the game itself, pausing it
func (p *Player) Tick(event tl.Event) {
	if event.Type == tl.EventKey { // Is it a keyboard event?
		switch event.Key { // If so, switch on the pressed key.
		case tl.KeyArrowRight:
			p.player.Current().Right()
		case tl.KeyArrowLeft:
			p.player.Current().Left()
		case tl.KeyArrowDown:
			p.player.Current().Down()
		case tl.KeyTab:
			p.player.Current().Rotate()
		case tl.KeySpace:
			if p.player.IsGameOver() {
				p.player.Reset()
				p.startGame()
			}
		}

		switch event.Ch {
		case 'p', 'P':
			p.player.Pause()
		}
	}
}
