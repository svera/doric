package main

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric/pkg/columns"
)

type Player struct {
	*tl.Entity
	player *columns.Player
}

func NewPlayer(p *columns.Player) *Player {
	return &Player{
		Entity: tl.NewEntity(0, 0, 0, 0),
		player: p,
	}
}

func (p *Player) Tick(event tl.Event) {
	if event.Type == tl.EventKey { // Is it a keyboard event?
		switch event.Ch { // If so, switch on the pressed key.
		case 'p':
		case 'P':
			p.player.Pause()
		}
	}
}
