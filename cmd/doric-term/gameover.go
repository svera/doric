package main

import (
	tl "github.com/JoelOtter/termloop"
)

type GameOver struct {
	*tl.Text
	game    *tl.Game
	restart func()
}

func NewGameOver(restart func(), offsetX int, offsetY int) *GameOver {
	return &GameOver{
		Text:    tl.NewText(offsetX, offsetY, "GAME OVER", tl.ColorWhite, tl.ColorBlack),
		game:    game,
		restart: restart,
	}
}

func (g *GameOver) Tick(event tl.Event) {
	if event.Type == tl.EventKey { // Is it a keyboard event?
		switch event.Key { // If so, switch on the pressed key.
		case tl.KeySpace:
			g.restart()
		}
	}
}
