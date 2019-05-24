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
var events chan *columns.Event
var game *tl.Game
var mainLevel *tl.BaseLevel
var gameOverLevel *tl.BaseLevel
var pitEntity *Pit
var pieceEntity *Piece
var nextPieceEntity *Next
var playerEntity *Player
var score *tl.Text

func main() {
	game = tl.NewGame()
	game.Screen().SetFps(60)
	pit := columns.NewPit(13, 6)
	player = columns.NewPlayer(pit)
	events = make(chan *columns.Event)
	score = tl.NewText(offsetX+10, offsetY, fmt.Sprintf("%d", player.Score()), tl.ColorWhite, tl.ColorBlack)
	setUpMainLevel()
	setUpGameOverScreen(pitEntity, nextPieceEntity)
	startGame()
	game.Start()
}

func setUpMainLevel() {
	mainLevel = tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
	})
	pitEntity = NewPit(player.Pit(), offsetX, offsetY)
	pieceEntity = NewPiece(player.Current(), offsetX, offsetY)
	nextPieceEntity = NewNext(player.Next(), offsetX+10, offsetY+4)
	playerEntity = NewPlayer(player)
	mainLevel.AddEntity(pitEntity)
	mainLevel.AddEntity(pieceEntity)
	mainLevel.AddEntity(nextPieceEntity)
	mainLevel.AddEntity(playerEntity)
	mainLevel.AddEntity(score)
}

func setUpGameOverScreen(pitEntity *Pit, nextPieceEntity *Next) {
	gameOverEntity := NewGameOver(
		func() {
			startGame()
			score.SetText(fmt.Sprintf("%d", 0))
		},
		offsetX+5,
		offsetY+5,
	)
	gameOverLevel = tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
	})
	gameOverLevel.AddEntity(pitEntity)
	gameOverLevel.AddEntity(nextPieceEntity)
	gameOverLevel.AddEntity(score)
	gameOverLevel.AddEntity(gameOverEntity)
}

func startGame() {
	game.Screen().SetLevel(mainLevel)
	player.Play(events)
	go func() {
		for {
			select {
			case ev := <-events:
				if ev.Name() == columns.Finished {
					game.Screen().SetLevel(gameOverLevel)
					return
				}
				if ev.Name() == columns.Scored {
					score.SetText(fmt.Sprintf("%d", player.Score()))
				}
			}
		}
	}()
}
