package main

import (
	"fmt"
	"math/rand"
	"time"

	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric/pkg/columns"
)

const (
	offsetX   = 32
	offsetY   = 5
	pithWidth = 6
	pitHeight = 13
)

var game *columns.Game
var mainLevel *tl.BaseLevel
var score *tl.Text
var level *tl.Text

func main() {
	cfg := columns.Config{
		PointsPerTile:           10,
		NumberTilesForNextLevel: 10,
		InitialSlowdown:         10,
		Frequency:               200 * time.Millisecond,
	}

	actions := make(chan int)
	app := tl.NewGame()
	app.Screen().SetFps(60)
	pit := columns.NewPit(pitHeight, pithWidth)
	current := columns.NewPiece(pit)
	next := columns.NewPiece(pit)
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	game = columns.NewGame(pit, *current, *next, r, cfg)
	score = tl.NewText(offsetX+15, offsetY, fmt.Sprintf("Score: %d", 0), tl.ColorWhite, tl.ColorBlack)
	level = tl.NewText(offsetX+15, offsetY+1, fmt.Sprintf("Level: %d", 1), tl.ColorWhite, tl.ColorBlack)
	pitEntity := NewPit(pit, offsetX, offsetY)
	message := tl.NewText(offsetX+1, offsetY+5, "", tl.ColorBlack, tl.ColorWhite)
	playerEntity := NewPlayer(current, actions, message, offsetX, offsetY)
	nextPieceEntity := NewNext(next, offsetX+15, offsetY+5)

	setUpMainLevel(pitEntity, playerEntity, nextPieceEntity, message)
	app.Screen().SetLevel(mainLevel)
	startGameLogic(actions, pitEntity, playerEntity, nextPieceEntity, message)
	app.Start()
}

func setUpMainLevel(pitEntity *Pit, playerEntity *Player, nextPieceEntity *Next, message *tl.Text) {
	mainLevel = tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
	})
	mainLevel.AddEntity(pitEntity)
	mainLevel.AddEntity(playerEntity)
	mainLevel.AddEntity(nextPieceEntity)
	mainLevel.AddEntity(score)
	mainLevel.AddEntity(level)
	mainLevel.AddEntity(message)
}

func startGameLogic(actions chan int, pitEntity *Pit, playerEntity *Player, nextPieceEntity *Next, message *tl.Text) {
	score.SetText(fmt.Sprintf("Score: %d", 0))
	level.SetText(fmt.Sprintf("Level: %d", 1))
	updates := make(chan columns.Update)
	go game.Play(actions, updates)

	go func() {
		defer func() {
			close(actions)
		}()
		for {
			select {
			case ev := <-updates:
				if ev.Status == columns.StatusFinished {
					playerEntity.Status = columns.StatusFinished
					return
				}
				if ev.Status == columns.StatusScored {
					score.SetText(fmt.Sprintf("Score: %d", ev.Points))
					level.SetText(fmt.Sprintf("Level: %d", ev.Level))
					pitEntity.Pit = ev.Pit
					playerEntity.Status = ev.Status
				}
				if ev.Status == columns.StatusPaused {
					playerEntity.Status = columns.StatusPaused
				}
				if ev.Status == columns.StatusUpdated {
					pitEntity.Pit = ev.Pit
					playerEntity.Current = &ev.Current
					playerEntity.Status = ev.Status
				}
				if ev.Status == columns.StatusRenewed {
					playerEntity.Current = &ev.Current
					nextPieceEntity.Piece = &ev.Next
				}
			}
		}
	}()
}
