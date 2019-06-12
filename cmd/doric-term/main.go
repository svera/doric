package main

import (
	"fmt"

	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric/pkg/columns"
)

const (
	offsetX   = 32
	offsetY   = 5
	pithWidth = 6
	pitHeight = 13
)

var player *columns.Player
var events chan int
var mainLevel *tl.BaseLevel
var score *tl.Text
var level *tl.Text

func main() {
	game := tl.NewGame()
	game.Screen().SetFps(60)
	pit := columns.NewPit(pitHeight, pithWidth)
	player = columns.NewPlayer(pit)
	events = make(chan int)
	score = tl.NewText(offsetX+15, offsetY, fmt.Sprintf("Score: %d", player.Score()), tl.ColorWhite, tl.ColorBlack)
	level = tl.NewText(offsetX+15, offsetY+1, fmt.Sprintf("Level: %d", player.Level()), tl.ColorWhite, tl.ColorBlack)
	setUpMainLevel()
	game.Screen().SetLevel(mainLevel)
	startGameLogic()
	game.Start()
}

func setUpMainLevel() {
	mainLevel = tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
	})
	pitEntity := NewPit(player.Pit(), offsetX, offsetY)
	message := tl.NewText(offsetX+1, offsetY+5, "", tl.ColorBlack, tl.ColorWhite)
	playerEntity := NewPlayer(player, startGameLogic, message, offsetX, offsetY)
	nextPieceEntity := NewNext(player.Next(), offsetX+15, offsetY+5)
	mainLevel.AddEntity(pitEntity)
	mainLevel.AddEntity(playerEntity)
	mainLevel.AddEntity(nextPieceEntity)
	mainLevel.AddEntity(score)
	mainLevel.AddEntity(level)
	mainLevel.AddEntity(message)
}

func startGameLogic() {
	score.SetText(fmt.Sprintf("Score: %d", player.Score()))
	level.SetText(fmt.Sprintf("Level: %d", player.Level()))
	go player.Play(events)
	go func() {
		for {
			select {
			case ev := <-events:
				if ev == columns.Finished {
					return
				}
				if ev == columns.Scored {
					score.SetText(fmt.Sprintf("Score: %d", player.Score()))
					level.SetText(fmt.Sprintf("Level: %d", player.Level()))
				}
			}
		}
	}()
}
