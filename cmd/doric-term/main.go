package main

import (
	"fmt"

	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric/pkg/columns"
)

const (
	offsetX = 35
	offsetY = 5
)

var player *columns.Player
var events chan int
var game *tl.Game
var mainLevel *tl.BaseLevel
var score *tl.Text
var level *tl.Text

func main() {
	game = tl.NewGame()
	game.Screen().SetFps(60)
	pit := columns.NewPit(13, 6)
	player = columns.NewPlayer(pit)
	events = make(chan int)
	score = tl.NewText(offsetX+10, offsetY, fmt.Sprintf("Score: %d", player.Score()), tl.ColorWhite, tl.ColorBlack)
	level = tl.NewText(offsetX+10, offsetY+1, fmt.Sprintf("Level: %d", player.Level()), tl.ColorWhite, tl.ColorBlack)
	setUpMainLevel()
	game.Screen().SetLevel(mainLevel)
	startGame()
	game.Start()
}

func setUpMainLevel() {
	mainLevel = tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
	})
	pitEntity := NewPit(player.Pit(), offsetX, offsetY)
	message := tl.NewText(offsetX+1, offsetY+5, "", tl.ColorBlack, tl.ColorWhite)
	playerEntity := NewPlayer(player, startGame, message, offsetX, offsetY)
	nextPieceEntity := NewNext(player.Next(), offsetX+10, offsetY+5)
	mainLevel.AddEntity(pitEntity)
	mainLevel.AddEntity(playerEntity)
	mainLevel.AddEntity(nextPieceEntity)
	mainLevel.AddEntity(score)
	mainLevel.AddEntity(level)
	mainLevel.AddEntity(message)
}

func startGame() {
	player.Play(events)
	score.SetText(fmt.Sprintf("Score: %d", player.Score()))
	level.SetText(fmt.Sprintf("Level: %d", player.Level()))
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
