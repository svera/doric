package main

import (
	"fmt"

	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric/pkg/columns"
)

const (
	offsetX = 10
	offsetY = 0
)

func main() {
	game := tl.NewGame()
	game.Screen().SetFps(60)
	mainLevel := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
	})
	pit := columns.NewPit(13, 6)
	player := columns.NewPlayer(pit)
	events := make(chan *columns.Event)
	player.Play(events)
	pitEntity := NewPit(player.Pit(), offsetX, offsetY)
	pieceEntity := NewPiece(player.Current(), offsetX, offsetY)
	nextPieceEntity := NewNext(player.Next(), offsetX+10, offsetY+4)
	score := tl.NewText(offsetX+10, offsetY, fmt.Sprintf("%d", player.Score()), tl.ColorWhite, tl.ColorBlack)
	mainLevel.AddEntity(pitEntity)
	mainLevel.AddEntity(pieceEntity)
	mainLevel.AddEntity(nextPieceEntity)
	mainLevel.AddEntity(score)
	game.Screen().SetLevel(mainLevel)
	go func() {
		for {
			select {
			case ev := <-events:
				if ev.Name() == columns.Finished {
					endLevel := tl.NewBaseLevel(tl.Cell{
						Bg: tl.ColorRed,
					})

					game.Screen().SetLevel(endLevel)
				}
				if ev.Name() == columns.Scored {
					score.SetText(fmt.Sprintf("%d", player.Score()))
				}
			}
		}
	}()
	game.Start()
}
